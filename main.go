package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultThreadNum = 10
	defaultRPS       = 2
)

type requestConfig struct {
	host      string
	threadNum int
	rps       int
}

func parseFlag(req *requestConfig) {
	flag.StringVar(&(req.host), "target", "http://127.0.0.1:8080/index.html", "target url")
	flag.IntVar(&(req.threadNum), "threadNum", defaultThreadNum, "number of threads to send requests")
	flag.IntVar(&req.rps, "rps", defaultRPS, "expected Requests-Per-Second")
	//flag.Set("logtostderr", "true")

	flag.Parse()
	glog.V(2).Infof("host=%s, rps=%d, threadNum=%d", req.host, req.rps, req.threadNum)
}

func handleSignal(pool *HttpReqPool, stop chan struct{}) {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		shouldForceStop := false

		for {
			sig := <-signalChannel
			switch sig {
			case os.Interrupt:
				glog.V(1).Info("Received SIGINT signal.")
			case syscall.SIGTERM:
				glog.V(1).Info("Reived SIGTERM signal, closed websocket")
			case syscall.SIGQUIT:
				glog.V(1).Info("Reived SIGQUIT signal, closed websocket")
			}
			if !shouldForceStop {
				pool.Stop()
				close(stop)
				shouldForceStop = true
			} else {
				glog.V(1).Info("Forced exit")
				fmt.Println("Forced exit")
				os.Exit(1)
			}
		}
	}()
}

func generateLoad(rps int, pool *HttpReqPool, stop chan struct{}) {
	ticker := time.NewTicker(time.Second * 1)

	i := int64(0)
	for {
		select {
		case <-ticker.C:
			end := i + int64(rps)
			for j := i; j < end; j++ {
				pool.Queue(j)
			}
			i = end
		case <-stop:
			glog.V(1).Infof("load generater is stopped.")
			return
		}
	}
}

func run(req *requestConfig) {
	hclient, err := NewHttpClient(req.host)
	if err != nil {
		glog.Errorf("Failed to create a http client: %v", err)
		return
	}

	f := func() {
		hclient.DoPost()
	}

	resp, err := hclient.DoPost()
	if err != nil {
		glog.Errorf("Error: %v", err)
		return
	} else {
		glog.V(2).Infof("response: %v", resp)
	}

	pool := NewHttpReqPool(req.threadNum, f)

	pool.Start()
	stop := make(chan struct{})
	handleSignal(pool, stop)
	generateLoad(req.rps, pool, stop)
	select {}
}

func main() {
	req := &requestConfig{}
	parseFlag(req)
	run(req)
}
