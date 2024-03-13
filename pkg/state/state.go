package state

import (
	"context"
	"errors"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/wrappers"
	"github.com/goccy/go-json"
)

// GetState is Naive Wrapper around Dapr State API
func GetState_deprecated(ctx context.Context, key string) []byte {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	// get store name from ctx
	store := ctx.Value("store").(string)
	item, err := client.GetState(ctx, store, key, nil)
	if err != nil {
		panic(err)
	}
	return item.Value
}

func GetState[T interface{}](ctx context.Context, key string) (T, error) {
	if common.CMEnabled {
		wrappers.PreRead(ctx, cm.Key(key))
	}
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	// get store name from ctx
	//store := ctx.Value("store").(string)
	item, err := client.GetState(ctx, common.RedisUrl, key, nil)
	if err != nil {
		panic(err)
	}
	var value T
	if len(item.Value) == 0 {
		return value, errors.New("key not found")
	}
	err = json.Unmarshal(item.Value, &value)
	if err != nil {
		panic(err)
	}
	return value, nil
}

func GetBulkState[T interface{}](ctx context.Context, keys []string) ([]T, error) {
	if common.CMEnabled {
		for _, key := range keys {
			wrappers.PreRead(ctx, cm.Key(key))
		}
	}
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	items, err := client.GetBulkState(ctx, common.RedisUrl, keys, nil, 10)
	if err != nil {
		panic(err)
	}
	values := make(map[string]*T)
	for _, item := range items {
		var value T
		if len(item.Value) == 0 {
			return nil, errors.New(fmt.Sprintf("key %s not found", item.Key))
		}
		err = json.Unmarshal(item.Value, &value)
		if err != nil {
			panic(err)
		}
		values[item.Key] = &value
	}
	var returnValues []T
	for _, key := range keys {
		returnValues = append(returnValues, *values[key])
	}
	return returnValues, nil
}

// This is a copy of the one above but without panicing for empty keys
func GetBulkStateDefault[T interface{}](ctx context.Context, keys []string, defVal T) []T {
	if common.CMEnabled {
		for _, key := range keys {
			wrappers.PreRead(ctx, cm.Key(key))
		}
	}
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	items, err := client.GetBulkState(ctx, common.RedisUrl, keys, nil, 10)
	if err != nil {
		panic(err)
	}
	values := make(map[string]*T)
	for _, item := range items {
		var value T
		if len(item.Value) == 0 {
			value = defVal
		} else {
			err = json.Unmarshal(item.Value, &value)
			if err != nil {
				panic(err)
			}
		}
		values[item.Key] = &value
	}
	var returnValues []T
	for _, key := range keys {
		returnValues = append(returnValues, *values[key])
	}
	return returnValues
}

func SetState(ctx context.Context, key string, value interface{}) {
	if common.CMEnabled {
		// prewrite is not necessary
		// wrappers.PreWrite(ctx, cm.Key(key))
	}
	valueBytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	//store := ctx.Value("store").(string)
	err = client.SaveState(ctx, common.RedisUrl, key, valueBytes, nil)
	if err != nil {
		panic(err)
	}
	if common.CMEnabled {
		wrappers.PostWrite(ctx, cm.Key(key))
	}
}

func SetBulkState(ctx context.Context, kvs map[string]interface{}) {
	items := make([]*dapr.SetStateItem, len(kvs))
	i := 0
	for k, v := range kvs {
		valueBytes, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		items[i] = &dapr.SetStateItem{
			Key:   k,
			Value: valueBytes,
		}
		i++
	}
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	err = client.SaveBulkState(ctx, common.RedisUrl, items...)
	if err != nil {
		panic(err)
	}
	if common.CMEnabled {
		for k, _ := range kvs {
			wrappers.PostWrite(ctx, cm.Key(k))
		}
	}
}
