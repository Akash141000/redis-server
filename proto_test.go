package main

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
	raw := []byte("*3\r\n$3\r\nSET\r\n$7\r\ntestKey\r\n$9\r\ntestValue\r\n")
	cmd, err := ParseCommand(raw)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Printf("cmd %+v", cmd)
}
