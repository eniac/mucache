package twoserivces

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/pkg/invoke"
	"time"
)

func InvalidationExperiment(times int, timeoutSeconds int, readApp string, readMethod string, writeApp string, writeMethod string) {
	timeout := time.Second * time.Duration(timeoutSeconds)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for key := 0; key < times; key++ {
		// Write a first value to populate the backend
		req0 := WriteRequest{K: key, V: 0}
		invoke.Invoke[string](ctx, writeApp, writeMethod, req0)
		//fmt.Printf("# Saved %+v for key %+v\n", req0.V, req0.K)
		//time.Sleep(time.Millisecond * 200)

		// Warm up the cache and read the latest value
		req1 := ReadRequest{K: key}
		resp := invoke.InvokeMissVal[ReadResponse](ctx, readApp, readMethod, req1)
		originalVal := resp.V
		//fmt.Printf("# Retrieved %+v for key %+v\n", resp.V, req1.K)

		// Should be adequate to wait until the cache is populated
		time.Sleep(time.Second * 5)
		invoke.AssertInvokeHit[ReadResponse](ctx, readApp, readMethod, req1)
		//fmt.Printf("# Retrieved %+v for key %+v from cache\n", cached.V, req1.K)

		// Take two timestamps to approximate the write time
		preWriteTs := time.Now()
		req2 := WriteRequest{K: key, V: originalVal + 1}
		invoke.Invoke[string](ctx, writeApp, writeMethod, req2)
		postWriteTs := time.Now()
		invoke.PollUntilMiss(ctx, readApp, readMethod, req1)
		postInvalidateTs := time.Now()

		diffWriteTs := postWriteTs.Sub(preWriteTs)
		meanWriteTs := preWriteTs.Add(diffWriteTs / 2)
		elapsedTime := postInvalidateTs.Sub(meanWriteTs)
		fmt.Printf("Invalidation elapsed time (in nanoseconds): %+v --- (writeDiff: %+v)\n", elapsedTime.Nanoseconds(), diffWriteTs.Nanoseconds())
	}
}
