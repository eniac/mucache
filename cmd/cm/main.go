package main

import (
	"flag"
	. "github.com/eniac/mucache/internal/cm"
	"github.com/eniac/mucache/pkg/common"
	_ "net/http/pprof"
)

func main() {
	var Port = flag.Int("port", 80,
		"The port to listen on when in http connection mode.")
	var cmAddressesPath = flag.String("cm_adds", "./experiments/local_cm/twoservices.txt",
		"The file that contains addresses of cache managers of different services.")
	var printTimeFreq = flag.Int("print_time_freq", 1000,
		"Every how many processed events should the cache manager print the time (for throughput measurements).")
	flag.Parse()
	//go func() {
	//	glog.Info(http.ListenAndServe("localhost:9090", nil))
	//}()

	cfg := InitConfig(*Port, *cmAddressesPath, *printTimeFreq)
	defer cfg.Close()

	// Initialize the state and call a processor
	state := NewState()
	go Process(cfg, state)
	go HttpSender(cfg)

	if common.ZMQ {
		// http servers are only used between cache managers
		go ServeHttp(cfg)
		Serve0mq(cfg)
		// Alternative native go implementation which is a little faster
		//Serve0mq2(cfg)
	} else {
		ServeHttp(cfg)
	}
}
