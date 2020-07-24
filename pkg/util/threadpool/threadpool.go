package threadpool

import (
	uuid "github.com/satori/go.uuid"
	"sync"
)

type Job func()

type ThreadPool struct {
	pool      map[string]interface{}
	poolMutex sync.Mutex
	jobs      map[string]Job
	jobMutex  sync.Mutex
	limit     int
	jobQueue  chan Job
}

func NewThreadPool(limit int) *ThreadPool {
	return &ThreadPool{
		pool:     map[string]interface{}{},
		jobs:     map[string]Job{},
		limit:    limit,
		jobQueue: make(chan Job, 1),
	}
}

func (t *ThreadPool) Run() {
	go t.run()
}

func (t *ThreadPool) run() {
	for {
		select {
		case job := <-t.jobQueue:
			go job()
		}
	}

}

func (t *ThreadPool) AddJob(job Job) {
	id := uuid.NewV4().String()
	t.jobs[id] = job
	t.jobQueue <- job
}
