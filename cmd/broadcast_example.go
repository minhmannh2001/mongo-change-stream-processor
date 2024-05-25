package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/minhmannh2001/mongo-change-stream-processor/pkg/broadcast"
)

func rangeChannel(
	ctx context.Context,
	n int,
) <-chan int {
	valueStream := make(chan int)
	go func() {
		defer close(valueStream)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case valueStream <- i:
			}
		}
	}()
	return valueStream
}

func broadcast_example() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Generates a channel sending integers
	// From 0 to 9
	range10 := rangeChannel(ctx, 10)

	broadcaster := broadcast.NewBroadcastServer(ctx, range10)
	listener1 := broadcaster.Subscribe()
	listener2 := broadcaster.Subscribe()
	listener3 := broadcaster.Subscribe()

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := range listener1 {
			fmt.Printf("Listener 1: %v/10 \n", i+1)
		}
	}()
	go func() {
		defer wg.Done()
		for i := range listener2 {
			fmt.Printf("Listener 2: %v/10 \n", i+1)
		}
	}()
	go func() {
		defer wg.Done()
		for i := range listener3 {
			fmt.Printf("Listener 3: %v/10 \n", i+1)
		}
	}()
	broadcaster.CancelSubscription(listener1)
	wg.Wait()
}

// https://betterprogramming.pub/how-to-broadcast-messages-in-go-using-channels-b68f42bdf32e
