package client

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})
	i, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	fmt.Println("i", i)
	return nil
}
