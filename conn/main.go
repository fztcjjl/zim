package main

import (
	"github.com/fztcjjl/zim/conn/app"
	"github.com/fztcjjl/zim/conn/handler/tcp"
	"github.com/fztcjjl/zim/conn/handler/ws"
)

func main() {
	a := app.NewApp()

	if srv := a.GetTcpServer(); srv != nil {
		srv.RegisterHandler(&tcp.Codec{}, tcp.NewHandler(a))
	}

	if srv := a.GetWsServer(); srv != nil {
		srv.RegisterHandler(&ws.Codec{}, ws.NewHandler(a))
	}

	a.Run()
}
