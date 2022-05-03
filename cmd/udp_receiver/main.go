package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/periweather/forza5"
)

var (
	Address string
	Port    int
)

func init() {
	flag.StringVar(&Address, "address", "127.0.0.1", "UDP server address")
	flag.IntVar(&Port, "port", 5607, "UDP server port")
	flag.Parse()
}

func main() {
	address := fmt.Sprintf("%s:%d", Address, Port)
	fmt.Printf("Starting server on %s\n", address)

	ctx := context.Background()

	err := forza5.Server(ctx, address)
	if err != nil {
		panic(err)
	}
}
