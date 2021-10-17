package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "server connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalln("There should be 2 arguments: host and port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		log.Println("Failed to connect:", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		client.Send()
		stop()
	}()

	go func() {
		client.Receive()
		stop()
	}()

	<-ctx.Done()
}
