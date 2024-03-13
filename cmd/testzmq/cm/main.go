package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"time"
)

func main() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5550")
	subscriber.SetSubscribe("")

	//  0MQ is so fast, we need to wait a while...
	time.Sleep(time.Second)
	syncclient, _ := zmq.NewSocket(zmq.REQ)
	defer syncclient.Close()
	syncclient.Connect("tcp://localhost:5551")
	syncclient.Send("", 0)
	syncclient.Recv(0)

	num := 0
	for {
		msg, err := subscriber.Recv(0)
		if err != nil {
			panic(err)
		}
		num += 1
		fmt.Println("Received: ", msg, num)
	}
}
