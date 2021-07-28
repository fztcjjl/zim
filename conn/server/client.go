package server

import (
	"bytes"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/pkg/ztimer"
	"github.com/panjf2000/gnet"
)

type Client struct {
	Status    int
	TimerTask *ztimer.TimerTask
	ConnId    string
	Conn      gnet.Conn
	Uin       string
	Platform  string
	Server    string
}

func (c *Client) Write(data []byte) error {
	return c.Conn.AsyncWrite(data)
}

func (c *Client) WritePacket(p *protocol.Proto) error {
	buf := &bytes.Buffer{}
	p.Write(buf)
	return c.Conn.AsyncWrite(buf.Bytes())
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
