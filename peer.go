package main

import (
	"fmt"
	"net"

	"golang.org/x/exp/slog"
)

type Peer struct {
	conn  net.Conn
	msgch chan []byte
}

func NewPeer(conn net.Conn, msg chan []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgch: msg,
	}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			fmt.Println("error", err.Error())
			if err.Error() != "EOF" {
				slog.Error("peer", "error peer read", err, "remoteAddr", p.conn.RemoteAddr())
				return err
			}
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgch <- msgBuf
		fmt.Println("buf", msgBuf)
	}
}
