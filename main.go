package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Swiffy12/taskify/src/internals/app"
	"github.com/Swiffy12/taskify/src/internals/config"
)

func main() {
	config := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := app.NewServer(config, ctx)

	go func() {
		osCall := <-c
		fmt.Printf("system call: %v\n", osCall)
		server.Shutdown()
		cancel()
	}()

	server.Listen()
}
