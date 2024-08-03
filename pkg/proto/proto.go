package proto

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tidwall/resp"
	"golang.org/x/exp/slog"
)

const (
	CommandSet = "set"
	CommandGet = "get"
)

type Command interface{}

type SetCommand struct {
	Key, Value []byte
}

type GetCommand struct {
	Key []byte
}

func ParseCommand(rawMsg []byte) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(string(rawMsg)))
	for {
		rv, _, err := rd.ReadValue()
		if err != nil && err.Error() != "EOF" {
			slog.Error("command", "error reading value", err)
			return nil, err
		}
		if rv.Type() == resp.Array {
			for _, value := range rv.Array() {
				switch strings.ToLower(value.String()) {
				case CommandSet:
					if len(rv.Array()) != 3 {
						return nil, fmt.Errorf("invalid number arguments for SET command")
					}
					cmd := SetCommand{
						Key:   rv.Array()[1].Bytes(),
						Value: rv.Array()[2].Bytes(),
					}
					return cmd, nil
				case CommandGet:
					if len(rv.Array()) != 2 {
						return nil, fmt.Errorf("invalid number arguments for GET command")
					}
					cmd := GetCommand{
						Key: rv.Array()[1].Bytes(),
					}
					return cmd, nil
				}
			}
		} else {
			return nil, fmt.Errorf("invalid command type")
		}
	}
}
