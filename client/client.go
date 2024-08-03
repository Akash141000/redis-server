package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		addr: addr,
		conn: conn,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error {
	fmt.Println("Set", "key", key)
	if c.conn == nil {
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			return err
		}
		c.conn = conn
	}

	buf := bytes.Buffer{}
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})
	i, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	io.Copy(&buf, c.conn)
	fmt.Println("buf", i)
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	fmt.Println("Get", "key", key)

	buf := bytes.Buffer{}
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})
	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return "", err
	}
	b := make([]byte, 1024)
	n, err := c.conn.Read(b)
	if err != nil {
		return "", err
	}
	return string(b[:n]), nil
}
