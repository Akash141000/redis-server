package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"redis/client"
	"redis/pkg/server"

	"golang.org/x/exp/slog"
)

const (
	testingKey   = "greet"
	testingValue = "hello!"
	addr         = "localhost:3000"
)

func clientSet() {
	c := client.New(addr)
	if err := c.Set(context.Background(), testingKey, testingValue); err != nil {
		log.Fatal(err)
	}
}

func clientGet() {
	c := client.New(addr)
	val, err := c.Get(context.Background(), testingKey)
	if err != nil {
		slog.Error("client", "error fetching the value", err)
		log.Fatal(err)
	}
	fmt.Println("val", val)

}

func main() {
	go func() {
		s := server.New(server.WithListenAddr(":3000"))
		s.Start()
	}()

	time.Sleep(time.Second * 3)
	go clientSet()

	time.Sleep(time.Second * 3)
	go clientGet()

	//block the server from exiting
	select {}
}
