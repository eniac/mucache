package invoke

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/eniac/mucache/pkg/wrappers"
	"github.com/golang/glog"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var daprPort = os.Getenv("DAPR_HTTP_PORT")

type SaveArgs struct {
	Ca     cm.CallArgs
	RetVal cm.ReturnVal
}

// Use a dedicated channel for saving to the cache in UpperBound mode
var SavingQueue = make(chan SaveArgs, 20000)

const SavingBatchSize = 10

// TODO: Add the cache-check in invocations. If the call is read-only, then first search for a value for specific
//       inputs in the cache. And if that does not exist, then we can simply perform the call.

// There is a problem with diamond shaped microservice graphs, when calls are performed to the same service
// in sequence. The problem is that we are not able to delay the write only for some part of the call,
// since this would lead into the service seeing incoherent state. Seeing the latest value, and then seeing a prior
// value.
//
// TODO: Modify the protocol to carry around in the ctx all of the services that have been visited throughout a request
//       also sending it upstream, so that it can then avoid the cache, if a cache-line has in its subtree a specific
//       service that we have visited (avoiding caching for merge points in microservice graph).

func Invoke_deprecated[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	if utility.IsCallReadOnly(app, method) {
		// TODO: Check the cache if the call is read-only
	}
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        inputBytes,
	}
	out, err := client.InvokeMethodWithContent(ctx, app, method, "post", content)
	if err != nil {
		panic(err)
	}
	var value T
	err = json.Unmarshal(out, &value)
	if err != nil {
		panic(err)
	}
	return value
}

// CheckCache checks if cache hit, return (value in cache, if hit)
func CheckCache(ctx context.Context, app string, method string, buf []byte) (cm.ReturnVal, bool) {
	if common.CMEnabled || common.UpperBoundEnabled {
		if utility.IsCallReadOnly(app, method) {
			ca := cm.CallArgs(wrappers.HashCallArgs(app, method, buf))
			ret, hit := wrappers.PreCall(ctx, ca)
			if hit {
				//glog.Info("Cache hit")
				return ret, true
			} else {
				//glog.Info("Cache miss")
				return "", false
			}
		}
	}
	return "", false
}

func CheckCacheShard(ctx context.Context, app string, method string, buf []byte) (cm.ReturnVal, int, bool) {
	if common.CMEnabled || common.UpperBoundEnabled {
		if utility.IsCallReadOnly(app, method) {
			ca, shard := wrappers.HashCallArgsWithShard(app, method, buf)
			ret, hit := wrappers.PreCall(ctx, cm.CallArgs(ca))
			if hit {
				//glog.Info("Cache hit")
				return ret, 0, true
			} else {
				//glog.Info("Cache miss")
				//glog.Infof("Cache miss, ca: %s, shard: %d", ca, shard)
				return "", shard, false
			}
		}
	}
	shard, err := strconv.Atoi(common.ShardCount)
	if err != nil {
		panic(err)
	}
	shard = rand.Int()%shard + 1 // shard is in [1, shardCount]
	return "", shard, false
}

func init() {
	if common.UpperBoundEnabled {
		go OffloadSave()
	}
}

func OffloadSave() {
	mc := cm.GetOrCreateCacheClient()
	expirationTTLmsInt, err := strconv.Atoi(common.ExpirationTTLms)
	if err != nil {
		expirationTTLmsInt = 0
	}
	//buf := make([]SaveArgs, 0, SavingBatchSize)
	expiration := time.Millisecond * time.Duration(expirationTTLmsInt)
	pipe := mc.Pipeline()
	counter := 0
	for {
		if len(SavingQueue) > 10000 {
			glog.Infof("SavingQueue size: %d", len(SavingQueue))
		}
		saveArgs := <-SavingQueue
		pipe.Set(context.Background(), string(saveArgs.Ca), string(saveArgs.RetVal), expiration)
		counter++
		if counter == SavingBatchSize {
			_, err := pipe.Exec(context.Background())
			if err != nil {
				panic(err)
			}
			pipe = mc.Pipeline()
			counter = 0
		}
		//err = mc.Set(context.Background(), string(saveArgs.Ca), string(saveArgs.RetVal), expiration).Err()
		//if err != nil {
		//	panic(err)
		//}
	}
}

func upperboundParseJsonAndSaveToCache(
	ctx context.Context, r io.Reader, target interface{},
	app string, method string, argBytes []byte) {
	var buf bytes.Buffer
	// Duplicate the body and read from it
	tee := io.TeeReader(r, &buf)
	utility.ParseJson(tee, target)
	// Save the response bytes
	retVal := cm.ReturnVal(buf.Bytes())
	// Create the call arguments
	ca := cm.CallArgs(wrappers.HashCallArgs(app, method, argBytes))
	SavingQueue <- SaveArgs{Ca: ca, RetVal: retVal}
	//mc := cm.GetOrCreateCacheClient()
	//expirationTTLmsInt, err := strconv.Atoi(common.ExpirationTTLms)
	////glog.Infof("Expiration string:%+v expiration int %+v\n", common.ExpirationTTLms, expirationTTLmsInt)
	//if err != nil {
	//	expirationTTLmsInt = 0
	//}
	//expiration := time.Millisecond * time.Duration(expirationTTLmsInt)
	////glog.Infof("Setting: %+v -> %+v with expiration: %+v\n", ca, retVal, expiration)
	//// Save the return value to the cache
	//err = mc.Set(ctx, string(ca), string(retVal), expiration).Err()
	//if err != nil {
	//	panic(err)
	//}
}

// Saves the response to *res (also might save the result to cache if we are in upperbound baseline
func performRequest[T interface{}](ctx context.Context, req *http.Request, res *T, app string, method string, argBytes []byte) {
	resp, err := common.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	utility.Assert(resp.StatusCode == http.StatusOK)
	defer resp.Body.Close()
	if common.UpperBoundEnabled && utility.IsCallReadOnly(app, method) {
		// If we are in the upper bound baseline implementation,
		// the caller saves its own cache if the call is read-only
		upperboundParseJsonAndSaveToCache(ctx, resp.Body, res, app, method, argBytes)
	} else {
		utility.ParseJson(resp.Body, res)
	}
}

func Invoke[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	if common.ShardEnabled {
		return ShardInvoke[T](ctx, app, method, input)
	}
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	value, hit := CheckCache(ctx, app, method, buf)
	var res T
	if hit {
		err := json.Unmarshal([]byte(value), &res)
		if err != nil {
			panic(err)
		}
		return res
	}
	url := fmt.Sprintf("http://localhost:%s/%s", daprPort, method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	req.Header.Set("dapr-app-id", app)
	if common.CMEnabled {
		req.Header.Set("caller", common.MyName)
		req.Header.Set("method", method)
	}
	performRequest[T](ctx, req, &res, app, method, buf)
	return res
}

func ShardInvoke[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	utility.Assert(common.ShardEnabled)
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	value, shard, hit := CheckCacheShard(ctx, app, method, buf)
	var res T
	if hit {
		err := json.Unmarshal([]byte(value), &res)
		if err != nil {
			panic(err)
		}
		return res
	}
	url := fmt.Sprintf("http://localhost:%s/%s", daprPort, method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	req.Header.Set("dapr-app-id", fmt.Sprintf("%s%d", app, shard))
	if common.CMEnabled {
		req.Header.Set("caller", common.MyName)
		req.Header.Set("method", method)
	}
	performRequest[T](ctx, req, &res, app, method, buf)
	return res
}

func InvokeHit(ctx context.Context, app string, method string, input interface{}) {
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	utility.Assert(common.CMEnabled || common.UpperBoundEnabled)
	utility.Assert(utility.IsCallReadOnly(app, method))
	ca := cm.CallArgs(wrappers.HashCallArgs(app, method, buf))
	//glog.Info("Start PreCall")
	wrappers.PreCall(ctx, ca) // ignore the return value
	//glog.Info("End PreCall")
	return
}

// Used for cache experiments
func AssertInvokeHit[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	utility.Assert(common.CMEnabled)
	utility.Assert(utility.IsCallReadOnly(app, method))
	ca := cm.CallArgs(wrappers.HashCallArgs(app, method, buf))
	//glog.Info("Start PreCall")
	ret, hit := wrappers.PreCall(ctx, ca) // ignore the return value
	//glog.Info("End PreCall")
	utility.Assert(hit)
	var value T
	err = json.Unmarshal([]byte(ret), &value)
	if err != nil {
		panic(err)
	}
	return value
}

func PollUntilMiss(ctx context.Context, app string, method string, input interface{}) {
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	utility.Assert(common.CMEnabled)
	utility.Assert(utility.IsCallReadOnly(app, method))
	ca := cm.CallArgs(wrappers.HashCallArgs(app, method, buf))
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Context timeout while polling for (%+v,%+v,%+v)!\n", app, method, input)
			return
		default:
			//glog.Info("Start PreCall")
			_, hit := wrappers.PreCall(ctx, ca) // ignore the return value
			//glog.Info("-- still hit")
			if !hit {
				return
			}
		}
	}
	return
}

func InvokeMiss[T interface{}](ctx context.Context, app string, method string, input interface{}) {
	InvokeMissVal[T](ctx, app, method, input)
	return
}

func InvokeMissVal[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	//utility.Assert(common.CMEnabled)
	if common.CMEnabled || common.UpperBoundEnabled {
		utility.Assert(utility.IsCallReadOnly(app, method))
		ca := cm.CallArgs(wrappers.HashCallArgs(app, method, buf))
		wrappers.PreCall(ctx, ca) // ignore the return value
	}
	url := fmt.Sprintf("http://localhost:%s/%s", daprPort, method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	req.Header.Set("dapr-app-id", app)
	if common.CMEnabled {
		req.Header.Set("caller", common.MyName)
		req.Header.Set("method", method)
	}
	var res T
	performRequest[T](ctx, req, &res, app, method, buf)
	//resp, err := common.HTTPClient.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//var value T
	//defer resp.Body.Close()
	//utility.ParseJson(resp.Body, &value)
	return res
}

func Invoke_deprecated2[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	if utility.IsCallReadOnly(app, method) {
		// TODO: Check the cache if the call is read-only
	}
	var buf bytes.Buffer
	utility.DumpJson(input, &buf)
	port := os.Getenv("DAPR_HTTP_PORT")
	url := fmt.Sprintf("http://localhost:%s/%s", port, method)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("dapr-app-id", app)

	mucacheOn := ctx.Value("mucache").(string)
	if mucacheOn == "on" {
		parentRID := ctx.Value("RID").(string)
		req.Header.Set("parentRID", parentRID)
	}
	req.Header.Set("mucache", "on")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	var value T
	utility.ParseJson(resp.Body, &value)
	return value
}

// What we want in python code
// dapr_port = os.getenv("DAPR_HTTP_PORT", 3500)
// dapr_url = "http://localhost:{}/neworder".format(dapr_port)
// message = {"data": {"orderId": n}}
// response = requests.post(dapr_url, json=message, timeout=5, headers = {"dapr-app-id": "nodeapp"} )

// Reference: https://docs.dapr.io/developing-applications/building-blocks/service-invocation/howto-invoke-discover-services/
