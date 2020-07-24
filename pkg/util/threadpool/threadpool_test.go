package threadpool

import (
	"fmt"
	"testing"
)

func TestNewThreadPool(t *testing.T) {
	pool := NewThreadPool(5)
	for i := 0; i < 10; i++ {
		pool.AddJob(func() {
			fmt.Printf("hello %d", i)
		})
	}
	pool.Run()
	select {}
}
