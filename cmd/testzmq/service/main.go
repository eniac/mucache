package main

import (
	zmq "github.com/pebbe/zmq4"
)

func main() {
	ctx, _ := zmq.NewContext()
	defer ctx.Term()
	publisher, _ := ctx.NewSocket(zmq.PUB)
	defer publisher.Close()
	publisher.SetSndhwm(1100000)
	publisher.Bind("tcp://*:5550")

	syncservice, _ := ctx.NewSocket(zmq.REP)
	defer syncservice.Close()
	syncservice.Bind("tcp://*:5551")

	syncservice.Recv(0)
	syncservice.Send("", 0)

	for i := 0; i < 1000; i++ {
		publisher.Send("Hello", 0)
	}
}
