package server

import (
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/conn/server/websocket"
	"github.com/panjf2000/gnet"
	"time"
)

type WsServer struct {
	gnet.EventHandler
	addr  string
	codec gnet.ICodec
	srv   *Server
}

func NewWsServer(srv *Server, addr string) *WsServer {
	ws := new(WsServer)
	ws.addr = addr
	ws.codec = &WsCodec{}
	ws.srv = srv
	return ws
}

func (ws *WsServer) Start() error {
	return gnet.Serve(ws, ws.addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(ws.codec))
}

func (ws *WsServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Debugf("ws server is listening on %s (multi-cores: %t, loops: %d)",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (ws *WsServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Debug("WS OnOpened ...")
	client := &Client{
		Status: WsUpgrading,
		Conn:   c,
	}
	c.SetContext(client)

	ws.srv.OnOpen(client)
	return
}

func (ws *WsServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Debug("WS OnClose ...")

	client, ok := c.Context().(*Client)
	if !ok {
		return
	}

	ws.srv.OnClose(client)
	return
}

func (ws *WsServer) React(data []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	client, ok := c.Context().(*Client)
	if !ok {
		return
	}

	ws.srv.OnMessage(data, client)

	return
}

// ==================================== Codec ==============================================

type WsCodec struct {
}

func (w *WsCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	s, ok := c.Context().(*Client)
	if !ok {
		return nil, nil
	}

	if s.Status != Authed {
		return buf, nil
	}

	f := websocket.NewBinaryFrame(buf)
	out, _ := websocket.FrameToBytes(&f)

	return out, nil
}

func (w *WsCodec) Decode(c gnet.Conn) ([]byte, error) {
	client, ok := c.Context().(*Client)
	if !ok {
		return nil, nil
	}

	if client.Status == WsUpgrading {
		r, out, err := websocket.ReadRequest(c)
		if err != nil {
			if err == websocket.ErrShortPackaet {
				return nil, nil
			}
			return out, err
		}
		out, err = websocket.Upgrade(c, r)
		c.AsyncWrite(out)
		if err == nil {
			client.Status = AuthPending
		}

		return nil, err
	} else {
		header, err := websocket.ReadHeader(c)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		_, payload := c.ReadN(int(header.Length))
		if header.Masked {
			websocket.Cipher(payload, header.Mask, 0)
		}

		if header.OpCode.IsControl() {
			switch header.OpCode {
			case websocket.OpClose:
				log.Debug("OnClose ...")
				//c.Close()
			case websocket.OpPing:
				log.Debug("OnPing ...")
			case websocket.OpPong:
				log.Debug("OpPong ...")
			}

			c.ShiftN(int(header.Length))
			return nil, nil
		}

		c.ShiftN(int(header.Length))
		return payload, nil
	}
}
