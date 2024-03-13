package cm

import (
	"bytes"
	"fmt"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/goccy/go-json"
	"github.com/golang/glog"
	zmq "github.com/pebbe/zmq4"
	"net/http"
)

const HttpStartSuffix = "/start"
const HttpEndSuffix = "/end"
const HttpInvSuffix = "/inv"
const HttpInvCallsSuffix = "/invcalls"
const HttpSaveCallsSuffix = "/save"

const QueueSize = 20000

// WQ is work queue; Go channel is thread-safe
// it's used both at Server and Client side
// Server side: it's used to receive requests from clients
// Client side: it's used to as a proxy to send requests to server
var WQ = make(chan interface{}, QueueSize)

func SetupZmqConnection() *zmq.Socket {
	ctx, _ := zmq.NewContext()
	publisher, _ := ctx.NewSocket(zmq.PUB)
	publisher.SetSndhwm(1100000)
	publisher.Bind("tcp://*:5550")

	syncservice, _ := ctx.NewSocket(zmq.REP)
	defer syncservice.Close()
	syncservice.Bind("tcp://*:5551")

	syncservice.Recv(0)
	syncservice.Send("", 0)
	glog.Info("ZMQ Connected")

	return publisher
}

func sendRequestHttp(req interface{}, url string) {
	buf := new(bytes.Buffer)
	utility.DumpJson(req, buf)
	resp, err := common.HTTPClient.Post(url, "application/json", buf)
	if err != nil {
		panic(err)
	}

	//Need to close the response stream, once response is read.
	//Hence defer close. It will automatically take care of it.
	defer resp.Body.Close()

	//Check response code, if New user is created then read response.
	utility.Assert(resp.StatusCode == http.StatusOK)
}

func ZmqProxy() {
	zmqPublisher := SetupZmqConnection()
	for {
		if len(WQ) >= QueueSize/2 {
			glog.Warningf("Queue size: %v\n", len(WQ))
		}
		req := <-WQ
		// Note: There is a tradeoff here w.r.t. where to do the json encoding.
		//       Doing it outside is multi-threaded but adds to critical path,
		//         while doing it here is single threaded but keeps it off critical path
		b, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}

		// TODO: We might want to batch send things other than start call here (to improve on the receiving end)
		_, err = zmqPublisher.Send(string(b), 0)
		if err != nil {
			panic(err)
		}
	}
}

func SendRequestZmq(req interface{}, t string) {
	fullReq := make(map[string]interface{})
	fullReq["type"] = t
	fullReq["inner"] = req
	WQ <- fullReq
}

func SendStartRequestHttp(req *StartRequest, url string) {
	utility.Assert(!common.ZMQ)
	sendRequestHttp(req, fmt.Sprintf("%s%s", url, HttpStartSuffix))
}

func SendEndRequestHttp(req *EndRequest, url string) {
	utility.Assert(!common.ZMQ)
	sendRequestHttp(req, fmt.Sprintf("%s%s", url, HttpEndSuffix))
}

func SendInvRequestHttp(req *InvalidateKeyRequest, url string) {
	//utility.Assert(!common.ZMQ)
	sendRequestHttp(req, fmt.Sprintf("%s%s", url, HttpInvSuffix))
}

func SendInvCallsRequestHttp(req *InvalidateCallsRequest, url string) {
	sendRequestHttp(req, fmt.Sprintf("%s%s", url, HttpInvCallsSuffix))
}

func SendSaveCallsRequestHttp(req *SaveCallsRequest, url string) {
	sendRequestHttp(req, fmt.Sprintf("%s%s", url, HttpSaveCallsSuffix))
}
