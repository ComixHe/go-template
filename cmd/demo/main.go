// Package main is the entry point for the demo application.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ComixHe/go-template/internal/metadata"
)

func main() {
	verPtr := flag.Bool("v", false, "print version and exit")
	flag.Parse()

	if *verPtr {
		fmt.Printf("%s", metadata.String())
		os.Exit(0)
	}

	fmt.Printf("Starting %s...\n", metadata.AppName)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-done
		fmt.Printf("\nReceived signal: %v. Shutting down...\n", sig)
		os.Exit(0)
	}()

	select {}
}
