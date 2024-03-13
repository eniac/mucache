package utility

import (
	"bytes"
	"encoding/gob"
	"github.com/goccy/go-json"
	"io"
	"strings"
)

func ParseJson(r io.Reader, target interface{}) {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		panic(err)
	}
}

func DumpJson(source interface{}, target io.Writer) {
	err := json.NewEncoder(target).Encode(source)
	if err != nil {
		panic(err)
	}
}

// This method checks if a call is read-only
//
// TODO: Make this read the read-only configuration from an environment variable or a file
func IsCallReadOnly(app string, method string) bool {
	return strings.HasPrefix(method, "ro_")
}

func Assert(cond bool) {
	if !cond {
		panic("Assertion failed!")
	}
}

func GetRealSizeOf(v interface{}) float32 {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		panic(err)
	}
	return float32(b.Len()) / 1024.0 / 1024.0
}

//func SetupCtx(r *http.Request, store string) ctx.Context {
//	ctx := r.Context()
//	ctx = ctx.WithValue(ctx, "store", store)
//	mucacheOn := r.Header.Get("mucache")
//	ctx = ctx.WithValue(ctx, "mucache", mucacheOn)
//	rid := r.Header.Get("RID")
//	ctx = ctx.WithValue(ctx, "RID", rid)
//	return ctx
//}

// RedisUrl returns the URL of the redis node with the given node index
// should only be called once
//func RedisUrl(nodeIdx string) string {
//	Assert(nodeIdx != "")
//	return fmt.Sprintf("redis%s", nodeIdx)
//}
//
//func MemcachedUrl(nodeIdx string) string {
//	Assert(nodeIdx != "")
//	return
//}
//
//func CacheManagerUrl(nodeIdx string) string {
//	Assert(nodeIdx != "")
//	return fmt.Sprintf("cm%s", nodeIdx)
//}
