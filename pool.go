package main

import (
	"github.com/golang/glog"
	"sync"
	"time"
)

type ThreadPool interface {
	Queue(int) error
	Start()
	Stop()
}

type jobfun func()

type HttpReqPool struct {
	Size    int
	stop    chan struct{}
	job     jobfun
	jobs    chan int64
	wg      sync.WaitGroup
	started bool

	mux sync.Mutex
}

func NewHttpReqPool(size int, job jobfun) *HttpReqPool {
	if size < 1 || job == nil {
		glog.Errorf("Invalid parameters.")
		return nil
	}

	result := &HttpReqPool{
		Size:    size,
		job:     job,
		stop:    make(chan struct{}),
		jobs:    make(chan int64, size*2),
		started: false,
	}

	return result
}

func (p *HttpReqPool) Start() {
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.started {
		glog.Errorf("Already started")
		return
	}
	p.started = true

	for i := 0; i < p.Size; i++ {
		j := i
		go p.worker(j)
	}

	p.wg.Add(p.Size)
}

func (p *HttpReqPool) worker(id int) {
	defer p.wg.Done()
	glog.V(2).Infof("Worker %d started", id)

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case jid, ok := <-p.jobs:
			if !ok {
				glog.V(2).Infof("Worker %d is noticfied jobs channel is closed.", id)
				return
			}
			glog.V(3).Infof("Worker %d gets a job-%d", id, jid)
			p.job()
			glog.V(3).Infof("Worker %d finishs a job-%d", id, jid)
		case <-p.stop:
			glog.V(2).Infof("Worker-%d received stop signal.", id)
			return
		case <-ticker.C:
			glog.V(3).Infof("Worker %d is waiting for job", id)
		}
	}
}

func (p *HttpReqPool) Stop() {
	p.mux.Lock()
	defer p.mux.Unlock()
	if !p.started {
		return
	}
	p.started = false

	close(p.stop)
	close(p.jobs)

	p.wg.Wait()
}

func (p *HttpReqPool) Queue(i int64) error {
	p.jobs <- i
	return nil
}
