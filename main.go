package main

import (
	"context"
	"log"
	"time"

	"redis/client"

	"golang.org/x/exp/slog"
)

func main() {
	slog.Info("Server", "starting")
	go func() {
		s := NewServer(WithListenAddr(":3000"))
		s.Start()
	}()
	time.Sleep(time.Second * 1)

	c := client.New("localhost:3000")
	if err := c.Set(context.Background(), "testKey", "testValue"); err != nil {
		log.Fatal(err)
	}

	//block the server from exiting
	select {}
}
