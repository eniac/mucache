package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/loadcm"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"time"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func startExperiment(ctx context.Context, req *loadcm.StartExperiment) *string {
	// Start running the zmqfeeder
	fmt.Printf("Starting experiment for: %v readers with duration: %ds\n", req.Readers, req.DurationSecs)
	go feedCM(req.Readers, req.DurationSecs)
	resp := "OK"
	return &resp
}

func feedCM(readers int, duration int) {
	timeout := time.Second * time.Duration(duration)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keys := make([]cm.Key, readers)
	for i := 0; i < readers; i++ {
		callArgs := cm.CallArgs(fmt.Sprintf("call-%d", i))
		key := cm.Key(fmt.Sprintf("key-%d", i))
		keys[i] = key
		fmt.Printf("Starting reader: %v\n", i)
		go reader(ctx, callArgs, key)
	}

	fmt.Printf("Starting writer with keys: %v\n", keys)
	go writer(ctx, keys)

	// Stop the experiment after the given duration
	select {
	case <-time.After(timeout):
		fmt.Println("Experiment timed out, cancelling context!")
	}
}

func reader(ctx context.Context, callArgs cm.CallArgs, key cm.Key) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Reader for %v is done!\n", key)
			return
		default:
			singleRequest(callArgs, key)
		}
	}
}

func writer(ctx context.Context, keys []cm.Key) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Writer is done!\n")
			return
		default:
			if i == len(keys) {
				i = 0
			}
			key := keys[i]
			singleInvalidate(key)
		}
	}
}

func singleRequest(callArgs cm.CallArgs, key cm.Key) {
	keyDeps := []cm.Key{key}
	// Send the start of a request
	cm.SendRequestZmq(&cm.StartRequest{CallArgs: callArgs}, cm.TypeStartRequest)
	endReq := cm.EndRequest{CallArgs: callArgs, KeyDeps: keyDeps, CallDeps: []cm.CallArgs{}, Caller: "stub", ReturnVal: "OK"}
	// Send the end of the request
	cm.SendRequestZmq(&endReq, cm.TypeEndRequest)
}

func singleInvalidate(key cm.Key) {
	cm.SendRequestZmq(&cm.InvalidateKeyRequest{Key: key}, cm.TypeInvRequest)
}

func main() {
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/start_experiment", wrappers.NonROWrapper[loadcm.StartExperiment, string](startExperiment))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
