package main

import (
	"context"
	"fmt"
	"log"

	"redis/client"
	"redis/pkg/server"

	"golang.org/x/exp/slog"
)

const (
	testingKey   = "greet"
	testingValue = "hello!"
	addr         = "locahost:3000"
)

func clientSet() {
	fmt.Println("set command")
	c := client.New(addr)
	if err := c.Set(context.Background(), testingKey, testingValue); err != nil {
		log.Fatal(err)
	}
}

func clientGet() {
	fmt.Println("get command")
	c := client.New(addr)
	val, err := c.Get(context.Background(), testingKey)
	if err != nil {
		slog.Error("client", "error fetching the value", err)
		log.Fatal(err)
	}
	fmt.Println("val", val)

}

func main() {
	slog.Info("Server", "starting")
	go func() {
		s := server.New(server.WithListenAddr(":3000"))
		s.Start()
	}()

	go clientSet()

	go clientGet()

	//block the server from exiting
	select {}
}
