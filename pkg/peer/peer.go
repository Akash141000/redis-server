package peer

import (
	"net"

	"golang.org/x/exp/slog"
)

type Message struct {
	Data []byte
	Peer *Peer
}

type Peer struct {
	Conn  net.Conn
	Msgch chan Message
}

func (p *Peer) Send(msg []byte) error {
	if _, err := p.Conn.Write(msg); err != nil {
		slog.Error("peer", "error sending message", err, "remoteAddr", p.Conn.RemoteAddr())
		return err
	}
	return nil
}

func New(conn net.Conn, msgch chan Message) *Peer {
	return &Peer{
		Conn:  conn,
		Msgch: msgch,
	}
}

func (p *Peer) ReadLoop() error {
	readBuf := make([]byte, 1024)
	for {
		n, err := p.Conn.Read(readBuf)
		if err != nil {
			if err.Error() != "EOF" {
				slog.Error("peer", "error peer read", err, "remoteAddr", p.Conn.RemoteAddr())
				return err
			}
			continue
		}

		//create msgBug and copy the readBuf
		msgBuf := make([]byte, n)
		copy(msgBuf, readBuf[:n])

		//create new msg and add to msgch
		msg := Message{
			Data: msgBuf,
			Peer: p,
		}
		p.Msgch <- msg
	}
}
