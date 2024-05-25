package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	barrier_example()

	<-quit
}
