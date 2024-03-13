package cm

import (
	"context"
	"github.com/eniac/mucache/pkg/common"
	"github.com/redis/go-redis/v9"
)

////
//
// This module contains functions for interacting with the cache
//
////

// var MemcachedClient *memcache.Client
var CacheClient *redis.Client

// Returns false if the value does not exist and true together with the value if it does
//
//	func CacheGet(mc *memcache.Client, ca CallArgs) (ReturnVal, bool) {
//		item, err := mc.Get(ca.ToString())
//		if err == nil {
//			return ByteArrayToRetVal(item.Value), true
//		} else {
//			// TODO: Return a better default
//			return ReturnVal(""), false
//		}
//	}
func CacheGet(c *redis.Client, ca CallArgs) (ReturnVal, bool) {
	item, err := c.Get(context.Background(), string(ca)).Result()
	//fmt.Printf("Got from cache: %+v, %+v for ca: %+v\n", item, err, ca)
	switch {
	case err == redis.Nil:
		return "", false
	case err != nil:
		panic(err)
	case err == nil:
		return ReturnVal(item), true
	}
	panic("unreachable")
}

//func CacheGet2(mc *memcache.Client, ca CallArgs) (ReturnVal, bool) {
//	ch := make(chan ReturnVal, 1)
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
//	defer cancel()
//	go func(ctx context.Context, ch chan ReturnVal) {
//		item, err := mc.Get(ca.ToString())
//		if err == nil {
//			ch <- ByteArrayToRetVal(item.Value)
//		} else {
//			// TODO: Return a better default
//			ch <- ReturnVal("")
//		}
//	}(ctx, ch)
//	select {
//	case <-ctx.Done():
//		return "", false
//	case ret := <-ch:
//		return ret, true
//	}
//}

// TODO: Maybe we can do that in a goroutine (since we don't need to wait for this to happen)
//
//	func CacheSet(mc *memcache.Client, ca CallArgs, ret ReturnVal) {
//		mc.Set(&memcache.Item{Key: ca.ToString(), Value: ret.ToByteArray()})
//	}
func CacheSet(c *redis.Client, ca CallArgs, ret ReturnVal) {
	err := c.Set(context.Background(), string(ca), string(ret), 0).Err()
	if err != nil {
		panic(err)
	}
}

//func CacheRemove(mc *memcache.Client, ca CallArgs) {
//	CacheRemoveCalls(mc, []CallArgs{ca})
//}

// TODO: Maybe we can do that in a goroutine (since we don't need to wait for this to happen)
//
//	func CacheRemoveCalls(mc *memcache.Client, calls []CallArgs) {
//		// TODO: There must be a way to do all deletes at once. I think we have to use binary protocol
//		//       (https://github.com/memcached/memcached/issues/245#issuecomment-272334627)
//		//       which as far as I understand will pipeline the requests. However, if we send requests
//		//       to be done in a different goroutine that should already be an improvement.
//		for _, ca := range calls {
//			// We are not interested in its return value (if the key did not exist)
//			mc.Delete(ca.ToString())
//		}
//	}
func CacheRemoveCalls(c *redis.Client, calls []CallArgs) {
	callsStr := make([]string, len(calls))
	for i, ca := range calls {
		callsStr[i] = string(ca)
	}
	err := c.Del(context.Background(), callsStr...).Err()
	if err != nil {
		panic(err)
	}
}

//func CacheSaveCalls(mc *memcache.Client, callArgsList []CallArgs, returnVals []ReturnVal) {
//	for i := range callArgsList {
//		callArgs := callArgsList[i]
//		returnVal := returnVals[i]
//		CacheSet(mc, callArgs, returnVal)
//	}
//}

func CacheSaveCalls(c *redis.Client, callArgsList []CallArgs, returnVals []ReturnVal) {
	callMap := make(map[string]string, len(callArgsList))
	for i, ca := range callArgsList {
		callMap[string(ca)] = string(returnVals[i])
	}
	err := c.MSet(context.Background(), callMap).Err()
	if err != nil {
		panic(err)
	}
}

//func GetOrCreateCacheClient() *memcache.Client {
//	if MemcachedClient == nil {
//		MemcachedClient = memcache.New(common.MemcachedUrl)
//		MemcachedClient.MaxIdleConns = 200
//		MemcachedClient.Timeout = 5 * time.Millisecond
//	}
//	return MemcachedClient
//}

func GetOrCreateCacheClient() *redis.Client {
	if CacheClient == nil {
		CacheClient = redis.NewClient(&redis.Options{
			Addr: common.CachedUrl,
		})
		err := CacheClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	return CacheClient
}
