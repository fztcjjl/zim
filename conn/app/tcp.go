package app

import (
	"github.com/panjf2000/gnet"
	"time"
)

type TcpServer struct {
	handler gnet.EventHandler
	addr    string
	codec   gnet.ICodec
}

func NewTcpServer(a *App, addr string) *TcpServer {
	srv := new(TcpServer)
	srv.addr = addr
	return srv
}

func (s *TcpServer) RegisterHandler(codec gnet.ICodec, handler gnet.EventHandler) {
	s.handler = handler
	s.codec = codec
}

func (s *TcpServer) Start() error {
	return gnet.Serve(s.handler, s.addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(s.codec))
}
