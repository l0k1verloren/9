package chanwg

import (
	"testing"
	"time"
)

func TestWaitGroup(t *testing.T) {
	wg := New(10)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		println("added waitgroup worker")
		go func() {
			// <-time.After(time.Second / 5)
			println("worker done")
			wg.Done(1)
		}()
	}
	go func() {
		<-time.After(time.Second)
		println("if this is above the end of the test output you failed :D")
	}()
	println("waiting for worker to finish")
	wg.Wait()
	println("done")
}