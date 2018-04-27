package main

import (
	"fmt"
	"testing"
	"time"
)

func job() {
	fmt.Printf(" job: %v\n", time.Now())
	time.Sleep(1)
}

func TestHttpReqPool_Start(t *testing.T) {
	p := NewHttpReqPool(5, job)
	p.Start()
	for i := int64(0); i < 10; i++ {
		p.Queue(i)
	}

	time.Sleep(time.Second * 20)
	p.Stop()
}
