package peer

import (
	"fmt"
	"net"
	"testing"
)

func TestPeer(t *testing.T) {
	fmt.Println("test peer")
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("tcp dial error", err)
	}

	writeMsg := []byte("writeMsgTest")
	n, err := conn.Write(writeMsg)
	if err != nil {
		fmt.Println("connection write error")
	}
	fmt.Println("bytes written", n)
}
