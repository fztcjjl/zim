package ws

import (
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/conn/app"
	"github.com/fztcjjl/zim/conn/handler/ws/websocket"
	"github.com/panjf2000/gnet"
)

type Codec struct {
}

func (w *Codec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	if c.Context() == nil {
		return buf, nil
	}
	f := websocket.NewTextFrame(buf)
	out, _ := websocket.FrameToBytes(&f)

	return out, nil
}

func (w *Codec) Decode(c gnet.Conn) ([]byte, error) {
	if c.Context() == nil {
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
			c.SetContext(app.AuthPending)
		}

		return nil, err
	} else {
		//status := c.Context().(int)
		header, err := websocket.ReadHeader(c)
		if err != nil {
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
