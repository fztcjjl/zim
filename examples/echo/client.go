package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/conn/server"
	"github.com/golang/protobuf/proto"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	headerLen = 20
)

type Client struct {
	user string
	addr string

	conn        net.Conn
	inputBuffer *Buffer
	seq         uint32
	syncSeq     int64
	authed      chan bool
}

func NewClient(user, addr string) *Client {
	client := new(Client)
	client.user = user
	client.addr = addr
	client.inputBuffer = NewBuffer(64 * 1024)
	client.authed = make(chan bool)
	return client
}

func (c *Client) Start() (err error) {
	if c.conn, err = net.Dial("tcp", c.addr); err != nil {
		log.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	go c.receive()
	if err = c.auth(); err != nil {
		log.Println(err)
		return
	}

	authed := <-c.authed
	if !authed {
		log.Println("auth error")
		return
	}
	c.sync()

	go c.heartbeat()

	stdin := bufio.NewReader(os.Stdin)
	users := []string{"李四", "王五", "赵六"}
	var seq uint32
	for {
		line, _, err := stdin.ReadLine()
		if err != nil {
			log.Printf("read from console failed, err: %v\n", err)
			break
		}
		seq = seq + 1

		to := users[rand.Int()%3]
		m := protocol.SendReq{
			ConvType:   1, // 会话类型：单聊
			MsgType:    1, // 文本消息
			From:       c.user,
			To:         to,
			Content:    string(line),
			ClientTime: time.Now().Unix(),
			Extra:      "",
		}

		mb, _ := proto.Marshal(&m)

		p := server.Packet{
			HeaderLen:     headerLen,
			ClientVersion: 1,
			Cmd:           uint32(protocol.CmdId_Cmd_SendReq),
			Seq:           seq,
			BodyLen:       uint32(len(mb)),
			Body:          mb,
		}

		if err := c.Write(&p); err != nil {
			log.Println(err)
			break
		}
	}

	return
}

func (c *Client) auth() (err error) {
	c.seq = 1
	auth := protocol.AuthReq{
		Uin:      c.user,
		Platform: "iOS",
		Token:    "dummy token",
	}

	b, err := proto.Marshal(&auth)
	if err != nil {
		return
	}

	p := server.Packet{
		HeaderLen:     headerLen,
		ClientVersion: 1,
		Cmd:           uint32(protocol.CmdId_Cmd_AuthReq),
		Seq:           c.seq,
		BodyLen:       uint32(len(b)),
		Body:          b,
	}

	err = c.Write(&p)
	c.seq++

	return

}

func (c *Client) sync() (err error) {
	req := protocol.SyncMsgReq{
		Offset: c.syncSeq,
		Limit:  20,
	}

	b, err := proto.Marshal(&req)
	if err != nil {
		return
	}

	p := server.Packet{
		HeaderLen:     headerLen,
		ClientVersion: 1,
		Cmd:           uint32(protocol.CmdId_Cmd_SyncMsgReq),
		Seq:           c.seq,
		BodyLen:       uint32(len(b)),
		Body:          b,
	}

	err = c.Write(&p)
	c.seq++

	return
}

func (c *Client) receive() {
	for {
		n, err := c.inputBuffer.Read(c.conn)
		if n == 0 {
			log.Println("Connection Closed")
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		data := c.inputBuffer.Peek()

		for {
			if len(data) < headerLen {
				break
			}

			bodyLen := binary.BigEndian.Uint32(data[16:headerLen])
			if len(data) < int(headerLen+bodyLen) {
				break
			}

			p := server.Packet{}
			p.Read(data[:headerLen+bodyLen])
			c.inputBuffer.Retrieve(int(headerLen + bodyLen))

			if p.Cmd == uint32(protocol.CmdId_Cmd_AuthRsp) {
				m := protocol.AuthRsp{}
				if err := proto.Unmarshal(p.Body, &m); err != nil {
					return
				}

				log.Println("auth rsp ...")

				c.authed <- true
			}
			if p.Cmd == uint32(protocol.CmdId_Cmd_SendRsp) {
				m := protocol.SendRsp{}
				if err := proto.Unmarshal(p.Body, &m); err != nil {
					return
				}

				log.Println("send rsp ...")
			}
			if p.Cmd == uint32(protocol.CmdId_Cmd_PushMsg) {
				m := protocol.Msg{}
				if err := proto.Unmarshal(p.Body, &m); err != nil {
					return
				}

				log.Printf("recv a msg: %s\n", m.Content)
			}
			if p.Cmd == uint32(protocol.CmdId_Cmd_Noop) {
				log.Println("noop ...")
			}
			if p.Cmd == uint32(protocol.CmdId_Cmd_SyncMsgRsp) {
				log.Println("sync rsp ...")
				m := protocol.SyncMsgRsp{}
				if err := proto.Unmarshal(p.Body, &m); err != nil {
					log.Println("error", err)
					return
				}
				size := len(m.List)
				log.Printf("sync %d msgs\n", size)

				if size > 0 {
					c.syncSeq = m.List[size-1].Seq
					c.sync()
				}
			}

			data = data[headerLen+bodyLen:]
		}

	}
}

func (c *Client) heartbeat() {

}

func (c *Client) Write(p *server.Packet) (err error) {
	buf := &bytes.Buffer{}
	p.Write(buf)
	_, err = c.conn.Write([]byte(buf.Bytes()))
	return
}
