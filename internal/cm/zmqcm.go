package cm

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	zmq2 "github.com/go-zeromq/zmq4"
	"github.com/goccy/go-json"
	"github.com/golang/glog"
	zmq "github.com/pebbe/zmq4"
	"time"
)

func Serve0mq2(cfg *Config) {
	ctx := context.Background()
	subscriber := zmq2.NewSub(ctx)
	defer subscriber.Close()

	subscriber.Dial(fmt.Sprintf("tcp://%v:%d", cfg.serviceName, cfg.port))
	subscriber.SetOption(zmq2.OptionSubscribe, "")
	//fmt.Printf("Subscribed to: %v:%d\n", cfg.serviceName, cfg.port)

	// allow the connection to build
	time.Sleep(time.Second)

	syncclient := zmq2.NewReq(ctx)
	defer syncclient.Close()
	syncclient.Dial(fmt.Sprintf("tcp://%v:%d", cfg.serviceName, cfg.port+1))
	emptyMsg := zmq2.NewMsgString("")
	syncclient.Send(emptyMsg)
	syncclient.Recv()

	glog.Info("ZMQ connected")

	for {
		msg, err := subscriber.Recv()
		if err != nil {
			panic(err)
		}
		handleMsgBytes(msg.Bytes())
	}
}

func Serve0mq(cfg *Config) {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect(fmt.Sprintf("tcp://%v:%d", cfg.serviceName, cfg.port))
	subscriber.SetSubscribe("")

	//fmt.Printf("Subscribed to: %v:%d\n", cfg.serviceName, cfg.port)

	// allow the connection to build
	time.Sleep(time.Second)

	// sync using a REQ/REP socket
	syncclient, _ := zmq.NewSocket(zmq.REQ)
	defer syncclient.Close()
	syncclient.Connect(fmt.Sprintf("tcp://%v:%d", cfg.serviceName, cfg.port+1))
	syncclient.Send("", 0)
	syncclient.Recv(0)

	glog.Info("ZMQ connected")

	for {
		msg, err := subscriber.Recv(0)
		if err != nil {
			panic(err)
		}
		handleMsg(msg)
	}
}

func handleMsg(msg string) {
	handleMsgBytes([]byte(msg))
}

func handleMsgBytes(msg []byte) {
	var typedReq map[string]json.RawMessage
	err := json.Unmarshal(msg, &typedReq)
	if err != nil {
		panic(err)
	}
	var t string
	err = json.Unmarshal(typedReq["type"], &t)
	if err != nil {
		panic(err)
	}
	switch t {
	case cm.TypeStartRequest:
		var req cm.StartRequest
		err = json.Unmarshal(typedReq["inner"], &req)
		if err != nil {
			panic(err)
		}
		cm.WQ <- req
	case cm.TypeEndRequest:
		var req cm.EndRequest
		err = json.Unmarshal(typedReq["inner"], &req)
		if err != nil {
			panic(err)
		}
		cm.WQ <- req
	case cm.TypeInvRequest:
		var req cm.InvalidateKeyRequest
		err = json.Unmarshal(typedReq["inner"], &req)
		if err != nil {
			panic(err)
		}
		cm.WQ <- req
	case cm.TypeInvCallsRequest:
		var req cm.InvalidateCallsRequest
		err = json.Unmarshal(typedReq["inner"], &req)
		if err != nil {
			panic(err)
		}
		cm.WQ <- req
	case cm.TypeSaveCallsRequest:
		var req cm.SaveCallsRequest
		err = json.Unmarshal(typedReq["inner"], &req)
		if err != nil {
			panic(err)
		}
		cm.WQ <- req
	default:
		panic("Unknown request type: " + t)
	}
}
