package proto

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
	rawMsg := []byte("*3\r\n$3\r\nSET\r\n$7\r\ntestKey\r\n$9\r\ntestValue\r\n")
	cmd, err := ParseCommand(rawMsg)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Printf("cmd %+v", cmd)
}
