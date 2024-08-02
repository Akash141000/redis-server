package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
	"golang.org/x/exp/slog"
)

const CommandSet = "Set"

type Command interface{}

type SetCommand struct {
	key, value string
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
						key:   value.Array()[1].String(),
						value: value.Array()[2].String(),
					}
					fmt.Println("cmd", cmd)
					return cmd, nil
				}
			}
		} else {
			return nil, fmt.Errorf("invalid command type")
		}
	}
	return nil, nil
}
