package app

import (
	"github.com/panjf2000/gnet"
	"time"
)

type WsServer struct {
	handler gnet.EventHandler
	addr    string
	codec   gnet.ICodec
}

func NewWsServer(a *App, addr string) *WsServer {
	srv := new(WsServer)
	srv.addr = addr
	return srv
}

func (s *WsServer) RegisterHandler(codec gnet.ICodec, handler gnet.EventHandler) {
	s.handler = handler
	s.codec = codec
}

func (s *WsServer) Start() error {
	return gnet.Serve(s.handler, s.addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(s.codec))
}
