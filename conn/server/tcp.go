package server

import (
	"bytes"
	"encoding/binary"
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/pkg/errors"
	"github.com/panjf2000/gnet"
	"time"
)

type TcpServer struct {
	gnet.EventHandler
	addr  string
	codec gnet.ICodec
	srv   *Server
}

func NewTcpServer(srv *Server, addr string) *TcpServer {
	ts := new(TcpServer)
	ts.addr = addr
	ts.codec = &TcpCodec{}
	ts.srv = srv
	return ts
}

func (s *TcpServer) Start() error {
	return gnet.Serve(s, s.addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(s.codec))
}

func (s *TcpServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Debugf("tcp server is listening on %s (multi-cores: %t, loops: %d)",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (s *TcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Debug("TCP OnOpened ...")
	client := &Client{
		Status: AuthPending,
		Conn:   c,
	}
	c.SetContext(client)

	s.srv.OnOpen(client)

	return
}

func (s *TcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Debug("TCP OnClose ...")

	client, ok := c.Context().(*Client)
	if !ok {
		return
	}

	s.srv.OnClose(client)
	return
}

func (s *TcpServer) React(data []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	client, ok := c.Context().(*Client)
	if !ok {
		return
	}

	s.srv.OnMessage(data, client)

	return
}

// ==================================== Codec ==============================================

type TcpCodec struct {
}

func (_ *TcpCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

func (_ *TcpCodec) Decode(c gnet.Conn) ([]byte, error) {
	if size, header := c.ReadN(protocol.HeaderLen); size == protocol.HeaderLen {
		byteBuffer := bytes.NewBuffer(header)
		var p protocol.Proto
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.HeaderLen); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.ClientVersion); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.Cmd); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.Seq); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.BodyLen); err != nil {
			return nil, err
		}

		protocolLen := int(protocol.HeaderLen + p.BodyLen)
		if size, data := c.ReadN(protocolLen); size == protocolLen {
			c.ShiftN(protocolLen)
			log.Info(data)
			log.Info(len(data), size, protocolLen)
			return data, nil
		}
		return nil, errors.New("not enough payload data")
	}

	return nil, errors.New("not enough header data")
}
