package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/wrappers"
	"time"
)

const HttpUrl = "http://localhost:8080"

func httpClient() {
	ctx := context.Background()
	ctx1 := wrappers.InitReqCtx(ctx, cm.CallId(1), "rid1", "service1", true)
	// Start Request 1
	//wrappers.PreReqStart(ctx1)
	// Perform a read
	wrappers.PreRead(ctx1, "k1")

	ctx2 := wrappers.InitReqCtx(ctx, cm.CallId(2), "rid2", "service1", true)
	// Start Request 2
	//wrappers.PreReqStart(ctx2)
	// Perform a read
	wrappers.PreRead(ctx2, "k3")

	// Request 3 (performs write)
	ctx3 := wrappers.InitReqCtx(ctx, cm.CallId(3), "write1", "service1", false)
	//wrappers.PreReqStart(ctx3)
	wrappers.PreWrite(ctx3, "k1")
	// Perform the write
	wrappers.PostWrite(ctx3, "k1")
	wrappers.PreReqEnd(ctx3, "OK")

	// Req2: Perform a read
	wrappers.PreRead(ctx2, "k2")
	wrappers.PreReqEnd(ctx2, "ret2")

	// Req1: Perform a read
	wrappers.PreRead(ctx1, "k3")
	wrappers.PreReqEnd(ctx1, "ret1")

	// Request 4 (performs write)
	ctx4 := wrappers.InitReqCtx(ctx, cm.CallId(4), "write2", "service1", false)
	//wrappers.PreReqStart(ctx4)
	wrappers.PreWrite(ctx4, "k3")
	// Perform the write
	wrappers.PostWrite(ctx4, "k3")
	wrappers.PreReqEnd(ctx4, "OK")
}

func httpClientNoWrapper() {
	command1 := cm.StartRequest{"rid1"}
	cm.SendStartRequestHttp(&command1, HttpUrl)
	command2 := cm.StartRequest{"rid2"}
	cm.SendStartRequestHttp(&command2, HttpUrl)
	command3 := cm.InvalidateKeyRequest{"k1"} // "I k1\n"
	cm.SendInvRequestHttp(&command3, HttpUrl)
	command4 := cm.EndRequest{"rid2", "caller2", []cm.Key{cm.Key("k2"), cm.Key("k3")}, "ret2"} //"E rid2 caller2 k2 k3\n"
	cm.SendEndRequestHttp(&command4, HttpUrl)
	command5 := cm.EndRequest{"rid1", "caller1", []cm.Key{cm.Key("k1"), cm.Key("k3")}, "ret1"} //"E rid1 caller1 k1 k3\n"
	cm.SendEndRequestHttp(&command5, HttpUrl)
	command6 := cm.InvalidateKeyRequest{"k3"}
	cm.SendInvRequestHttp(&command6, HttpUrl)
}

func testMemcached() {
	mc := cm.GetOrCreateCacheClient()
	defer mc.Close()

	ca1 := cm.CallArgs("req1")
	ret1 := cm.ReturnVal("ret1")

	value, exists := cm.CacheGet(mc, ca1)
	fmt.Println("Get:", ca1, "Exists:", exists, "value:", value)
	cm.CacheSet(mc, ca1, ret1)
	fmt.Println("Set:", ca1, "to value:", ret1)
	value, exists = cm.CacheGet(mc, ca1)
	fmt.Println("Get:", ca1, "Exists:", exists, "value:", value)
}

// This is an experiment with two services (they don't really exist), and therefore 2 caches and 2 cache managers
func httpTwoService() {
	// Service 1 configuration
	bg1 := context.Background()

	// Service 2 configuration
	bg2 := context.Background()

	inS1 := func() {
		common.CMUrl = "http://localhost:8080"
		common.MemcachedUrl = "localhost:11211"
	}

	inS2 := func() {
		common.CMUrl = "http://localhost:8081"
		common.MemcachedUrl = "localhost:11212"
	}

	// ca11 performs three calls, to ca1, ca2, ca2, and ca2 again (the second ca2 cache-hits and the third should cache-miss)
	ca11 := cm.CallArgs("rid11")
	ca1 := cm.CallArgs("read_only_rid1")
	ca2 := cm.CallArgs("read_only_rid2")

	inS1()
	// This is the parent of call1
	ctx11 := wrappers.InitReqCtx(bg1, cm.CallId(11), ca11, "client", false)
	// Start Request 1
	//wrappers.PreReqStart(ctx11)
	// Make call rid1 to service 2 (should be a cache-miss)
	cacheV1, exists := wrappers.PreCall(ctx11, ca1)
	if exists {
		fmt.Println("Cache hit for:", ca1, "value:", cacheV1)
		panic("Cache should never hit here, make sure to restart memcached before running this experiment!")
	} else {
		fmt.Println("Cache miss for:", ca1)
	}

	inS2()

	ctx1 := wrappers.InitReqCtx(bg1, cm.CallId(1), ca1, "service1", true)
	// Start Request 1
	//wrappers.PreReqStart(ctx1)
	// Perform a read
	wrappers.PreRead(ctx1, "k1")

	inS1()
	// Make call rid2 to service 2 (should be a cache-miss)
	cacheV2, exists := wrappers.PreCall(ctx11, ca2)
	if exists {
		fmt.Println("Cache hit for:", ca2, "value:", cacheV2)
		panic("Cache should never hit here, make sure to restart memcached before running this experiment!")
	} else {
		fmt.Println("Cache miss for:", ca2)
	}

	inS2()
	ctx2 := wrappers.InitReqCtx(bg2, cm.CallId(2), ca2, "service1", true)
	// Start Request 2
	//wrappers.PreReqStart(ctx2)
	// Perform a read
	wrappers.PreRead(ctx2, "k3")

	// Request 3 (performs write)
	ctx3 := wrappers.InitReqCtx(bg1, cm.CallId(3), "write1", "service1", false)
	//wrappers.PreReqStart(ctx3)
	wrappers.PreWrite(ctx3, "k1")
	// Perform the write
	wrappers.PostWrite(ctx3, "k1")
	wrappers.PreReqEnd(ctx3, "OK")

	// Req2: Perform a read
	wrappers.PreRead(ctx2, "k2")
	wrappers.PreReqEnd(ctx2, "ret2")

	// Sleep for one second to make sure that the cache has been saved
	time.Sleep(500 * time.Millisecond)

	// Req1: Perform a read
	wrappers.PreRead(ctx1, "k3")
	wrappers.PreReqEnd(ctx1, "ret1")
	// This is invalid so it should not be saved

	inS1()
	// Make second call rid2 to service 2 (should be a cache-hit)
	cacheV2new, exists := wrappers.PreCall(ctx11, ca2)
	if exists {
		fmt.Println("Cache hit for:", ca2, "value:", cacheV2new)
	} else {
		fmt.Println("Cache miss for:", ca2)
		panic("Cache should never miss here!")
	}

	inS2()
	// Request 4 (performs write)
	ctx4 := wrappers.InitReqCtx(bg2, cm.CallId(4), "write2", "service1", false)
	//wrappers.PreReqStart(ctx4)
	wrappers.PreWrite(ctx4, "k3")
	// Perform the write
	wrappers.PostWrite(ctx4, "k3")
	wrappers.PreReqEnd(ctx4, "OK")
	time.Sleep(500 * time.Millisecond)

	inS1()
	// Make third call rid2 to service 2 (should be a cache-miss)
	cacheV2new3, exists := wrappers.PreCall(ctx11, ca2)
	if exists {
		fmt.Println("Cache hit for:", ca2, "value:", cacheV2new3)
		panic("Cache should never miss here!")
	} else {
		fmt.Println("Cache miss for:", ca2)
	}
}

func main() {
	// Parse arguments
	var connectionMode = flag.String("scenario", "http",
		"Define the execution scenario. Options: http_no_wrapper, http, test_memcached, http_two_service, unix_sock (not supported anymore).")
	flag.Parse()

	// Serve requests
	switch *connectionMode {
	case "http":
		httpClient()
	case "http_no_wrapper":
		httpClientNoWrapper()
	case "test_memcached":
		testMemcached()
	case "http_two_service":
		httpTwoService()
	default:
		panic(fmt.Sprintf("Unknown connection mode: %s", *connectionMode))
	}
}
