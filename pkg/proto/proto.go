package proto

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
	"golang.org/x/exp/slog"
)

const (
	CommandSet = "Set"
	CommandGet = "Get"
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
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("command", "error reading value", err)
			return nil, err
		}
		if rv.Type() == resp.Array {
			for _, value := range rv.Array() {
				switch value.String() {
				case CommandSet:
					if len(value.Array()) != 3 {
						return nil, fmt.Errorf("invalid number arguments for SET command")
					}
					cmd := SetCommand{
						Key:   value.Array()[1].Bytes(),
						Value: value.Array()[2].Bytes(),
					}
					return cmd, nil
				case CommandGet:
					if len(value.Array()) != 2 {
						return nil, fmt.Errorf("invalid number arguments for GET command")
					}
					cmd := GetCommand{
						Key: value.Array()[1].Bytes(),
					}
					return cmd, nil
				}
			}
		} else {
			return nil, fmt.Errorf("invalid command type")
		}
	}
	return nil, nil
}
