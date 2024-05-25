package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/minhmannh2001/mongo-change-stream-processor/pkg/barrier"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type counter struct {
	c int
	sync.Mutex
}

func (c *counter) Incr() {
	c.Lock()
	c.c += 1
	c.Unlock()
}
func (c *counter) Get() (res int) {
	c.Lock()
	res = c.c
	c.Unlock()
	return
}

func worker(c *counter, br *barrier.Barrier, wg *sync.WaitGroup) {
	for i := 0; i < 3; i++ {
		br.Before()
		c.Incr()
		br.After()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		fmt.Println(c.Get())
	}
	wg.Done()
}

func barrier_example() {
	var wg sync.WaitGroup
	workers := 3
	br := barrier.New(workers)
	c := counter{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(&c, br, &wg)
	}
	wg.Wait()
}

// https://medium.com/golangspec/reusable-barriers-in-golang-156db1f75d0b
