package wrappers

import (
	"context"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/lithammer/shortuuid"
	"io"
	"net/http"
)

//
// The following includes all of the ctx code
//
// KK: Let's contain all ctx code here so that we don't have to guess which key was used for which field.
//

const callArgsCtxIndex = "call-args"
const readOnlyCtxIndex = "read-only"
const callIdCtxIndex = "RID"
const callerCtxIndex = "caller"

type CtxItem struct {
	Key   string
	Value interface{}
}

func ReadOnlyContext(ctx context.Context) bool {
	readOnly := ctx.Value(readOnlyCtxIndex)
	if readOnly == nil {
		return false
	}
	return readOnly.(bool)
}

func CtxCallArgs(ctx context.Context) cm.CallArgs {
	// This can never return nil
	return ctx.Value(callArgsCtxIndex).(cm.CallArgs)
}

func CtxCallId(ctx context.Context) cm.CallId {
	// This can never return nil
	return ctx.Value(callIdCtxIndex).(cm.CallId)
}

func CtxCaller(ctx context.Context) cm.ServiceName {
	return ctx.Value(callerCtxIndex).(cm.ServiceName)
}

// use slice of pointer to avoid unnecessary copying
// see https://go.dev/play/p/OG4THyUq-lv
func setCtxItems(ctx context.Context, items []*CtxItem) context.Context {
	for _, item := range items {
		ctx = context.WithValue(ctx, item.Key, item.Value)
	}
	return ctx
}

// Helper function that initializes the ctx for the client
// It's used at the beginning of each request handler
// so PreReqStart is also merged into this function
func InitReqCtx(ctx context.Context, id cm.CallId, ca cm.CallArgs, caller cm.ServiceName, ro bool) context.Context {
	ctxItems := []*CtxItem{
		{callIdCtxIndex, id},
		{callArgsCtxIndex, ca},
		{callerCtxIndex, caller},
		{readOnlyCtxIndex, ro},
	}
	ctx = setCtxItems(ctx, ctxItems)
	PreReqStart(ctx)
	return ctx
}

func SetupCtxFromHTTPReq(r *http.Request, ro bool) (context.Context, []byte) {
	ctx := r.Context()
	inputBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	if common.CMEnabled {
		caller := r.Header.Get("caller")
		if caller == "" {
			caller = "client"
		}
		//id := shortuuid.NewWithNamespace(caller)
		id := shortuuid.New()
		method := r.Header.Get("method")
		var ca string
		if common.ShardEnabled {
			utility.Assert(common.MyRawName != "")
			ca = HashCallArgs(common.MyRawName, method, inputBytes)
		} else {
			ca = HashCallArgs(common.MyName, method, inputBytes)
		}
		ctx = InitReqCtx(ctx, cm.CallId(id), cm.CallArgs(ca), cm.ServiceName(caller), ro)
	}
	return ctx, inputBytes
}

// No need to use the manual methods

func CtxSetReadOnly(ctx context.Context, readOnly bool) context.Context {
	ctx = context.WithValue(ctx, readOnlyCtxIndex, readOnly)
	return ctx
}

func CtxSetCallArgs(ctx context.Context, ca cm.CallArgs) context.Context {
	ctx = context.WithValue(ctx, callArgsCtxIndex, ca)
	return ctx
}

func CtxSetCallId(ctx context.Context, id cm.CallId) context.Context {
	ctx = context.WithValue(ctx, callIdCtxIndex, id)
	return ctx
}

func CtxSetCaller(ctx context.Context, name cm.ServiceName) context.Context {
	ctx = context.WithValue(ctx, callerCtxIndex, name)
	return ctx
}
