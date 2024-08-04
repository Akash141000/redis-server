package peer

import (
	"fmt"
	"net"
	"redis/pkg/proto"
	"redis/pkg/store"
	"strings"

	"github.com/tidwall/resp"
	"golang.org/x/exp/slog"
)

type Peer struct {
	Conn  net.Conn
	store store.Storer
}

func (p *Peer) Send(msg []byte) error {
	if _, err := p.Conn.Write(msg); err != nil {
		slog.Error("peer", "error sending message", err, "remoteAddr", p.Conn.RemoteAddr())
		return err
	}
	return nil
}

func New(conn net.Conn, store store.Storer) *Peer {
	return &Peer{
		Conn:  conn,
		store: store,
	}
}

func (p *Peer) ReadLoop() error {
	rd := resp.NewReader(p.Conn)
	for {
		rv, _, err := rd.ReadValue()
		if err != nil && err.Error() != "EOF" {
			slog.Error("command", "error reading value", err)
			return err
		}
		if rv.Type() == resp.Array {
		rangeCommand:
			for _, value := range rv.Array() {
				switch strings.ToLower(value.String()) {
				case proto.CommandSet:
					if len(rv.Array()) != 3 {
						return fmt.Errorf("invalid number arguments for SET command")
					}
					cmd := proto.SetCommand{
						Key:   rv.Array()[1].Bytes(),
						Value: rv.Array()[2].Bytes(),
					}
					p.handleCommand(cmd)
					break rangeCommand
					// return nil
				case proto.CommandGet:
					if len(rv.Array()) != 2 {
						return fmt.Errorf("invalid number arguments for GET command")
					}
					cmd := proto.GetCommand{
						Key: rv.Array()[1].Bytes(),
					}
					p.handleCommand(cmd)
					break rangeCommand
				}
			}
		} else {
			return fmt.Errorf("invalid command type")
		}
	}
	//

}

func (p *Peer) handleCommand(cmd proto.Command) error {
	switch v := cmd.(type) {
	case proto.SetCommand:
		return p.store.Set(string(v.Key), v.Value)
	case proto.GetCommand:
		val, err := p.store.Get(v.Key)
		if err != nil {
			return err
		}
		//send the value found for the key over connection
		if err := p.Send(val); err != nil {
			return err
		}

	}

	return nil
}
