package tcp

import (
	"bytes"
	"context"
	"errors"
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/conn/app"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

type Handler struct {
	*gnet.EventServer
	app        *app.App
	workerPool *goroutine.Pool
}

func NewHandler(a *app.App) *Handler {
	h := new(Handler)
	h.app = a
	h.workerPool = goroutine.Default()
	return h
}

func (h *Handler) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Debugf("tcp server is listening on %s (multi-cores: %t, loops: %d)",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (h *Handler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Debug("TCP OnOpened ...")
	c.SetContext(app.AuthPending)
	return
}

func (h *Handler) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Debug("TCP OnClose ...")
	s := h.app.DelSessionByConn(c)

	status, ok := c.Context().(int)
	if !ok {
		return
	}

	if status != app.Authed {
		return
	}

	h.workerPool.Submit(func() {
		if s != nil {
			logicClient := h.app.GetLogicClient()
			req := logic.DisconnectReq{
				Uin:      s.Uin,
				Platform: s.Platform,
			}
			logicClient.Disconnect(context.Background(), &req)
		}
	})

	return
}

func (h *Handler) React(data []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	status, ok := c.Context().(int)
	if !ok {
		return
	}

	h.workerPool.Submit(func() {
		p := &protocol.Proto{}
		if err := p.Read(data); err != nil {
			c.Close()
			return
		}
		if status == app.AuthPending {
			if err := h.handleAuth(c, p); err != nil {
				c.Close()
				return
			} else {
				c.SetContext(app.Authed)
				return
			}
		} else {
			log.Info(p)
			h.handleProto(c, p)
		}
	})

	return
}

func (h *Handler) handleAuth(c gnet.Conn, p *protocol.Proto) (err error) {
	log.Info("handleAuth ...")
	req := &protocol.AuthReq{}

	rsp := &protocol.AuthRsp{
		Code:    0,
		Message: "成功",
	}

	defer func() {
		b, err := proto.Marshal(rsp)
		if err != nil {
			return
		}

		p.Cmd = uint32(protocol.CmdId_Cmd_AuthRsp)
		p.BodyLen = uint32(len(b))
		p.Body = b
		buf := &bytes.Buffer{}
		p.Write(buf)
		c.AsyncWrite(buf.Bytes())
	}()

	if err = proto.Unmarshal(p.Body, req); err != nil {
		log.Error(err)
		rsp.Code = -1
		rsp.Message = "协议解析错误"
		return
	}

	if req.Uin == "" {
		rsp.Code = -1
		rsp.Message = "账号不能为空"
		log.Error("账号不能为空")
		return errors.New("账号不能为空")
	}

	logicClient := h.app.GetLogicClient()
	reqL := logic.ConnectReq{
		Uin:      req.Uin,
		Platform: req.Platform,
		Server:   h.app.GetServerId(),
		Token:    req.Token,
		Device:   req.Device,
	}
	rspL, err := logicClient.Connect(context.Background(), &reqL)
	if err != nil {
		log.Error(err)
		rsp.Code = -1
		rsp.Message = "系统错误"
		return
	}

	log.Info("auth succ")
	rsp.Code = rspL.Code
	rsp.Message = rspL.Message
	// 踢掉旧的连接
	if rspL.KickedConnId != "" {
		s := h.app.DelSessionByConnId(rspL.KickedConnId)
		if s != nil && s.Conn != nil {
			kick := &protocol.Kick{KickReason: rspL.KickedReason}
			if b, err := proto.Marshal(kick); err != nil {
				s.Conn.AsyncWrite(b)
			}
			log.Info("close old")
			s.Conn.Close()
		}
	}

	s := &app.Session{
		ConnId:   rspL.ConnId,
		Conn:     c,
		Uin:      reqL.Uin,
		Platform: reqL.Platform,
		Server:   h.app.GetServerId(),
	}
	h.app.AddSession(s)

	return
}

func (h *Handler) handleProto(c gnet.Conn, p *protocol.Proto) (err error) {
	if p.Cmd == uint32(protocol.CmdId_Cmd_Noop) {
		err = h.handleNoop(c, p)
	} else if p.Cmd == uint32(protocol.CmdId_Cmd_SendReq) {
		err = h.handleSend(c, p)
	} else if p.Cmd == uint32(protocol.CmdId_Cmd_SyncMsgReq) {
		err = h.handleSyncMsg(c, p)
	} else if p.Cmd == uint32(protocol.CmdId_Cmd_MsgAckReq) {
		err = h.handleMsgAckReq(c, p)
	}

	return
}

func (h *Handler) handleMsgAckReq(c gnet.Conn, p *protocol.Proto) (err error) {
	return
}

func (h *Handler) handleSyncMsg(c gnet.Conn, p *protocol.Proto) (err error) {
	req := &protocol.SyncMsgReq{}
	rsp := &protocol.SyncMsgRsp{}

	defer func() {
		b, err := proto.Marshal(rsp)
		if err != nil {
			return
		}

		p.Cmd = uint32(protocol.CmdId_Cmd_SyncMsgRsp)
		p.BodyLen = uint32(len(b))
		p.Body = b
		buf := &bytes.Buffer{}
		p.Write(buf)
		c.AsyncWrite(buf.Bytes())
	}()

	if err = proto.Unmarshal(p.Body, req); err != nil {
		return
	}

	s := h.app.GetSessionByConn(c)
	logicClient := h.app.GetLogicClient()
	reqL := logic.SyncMsgReq{
		Uin:    s.Uin,
		ConnId: s.ConnId,
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	rspL, err := logicClient.SyncMsg(context.Background(), &reqL)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(rspL)

	for _, v := range rspL.List {
		msg := &protocol.Msg{
			Id:       v.Id,
			ConvType: v.ConvType,
			Type:     v.Type,
			Content:  v.Content,
			From:     v.From,
			To:       v.To,
			Extra:    v.Extra,
			SendTime: v.SendTime,
			Seq:      v.Seq,
		}
		rsp.List = append(rsp.List, msg)
	}

	return
}

func (h *Handler) handleNoop(c gnet.Conn, p *protocol.Proto) (err error) {
	buf := &bytes.Buffer{}
	p.Write(buf)
	c.AsyncWrite(buf.Bytes())

	s := h.app.GetSessionByConn(c)
	if s != nil {
		logicClient := h.app.GetLogicClient()
		req := logic.HeartbeatReq{
			Uin:    s.Uin,
			ConnId: s.ConnId,
			Server: s.Server,
		}
		logicClient.Heartbeat(context.Background(), &req)
	}

	return
}

func (h *Handler) handleSend(c gnet.Conn, p *protocol.Proto) (err error) {
	log.Info("handleSend ...")
	req := &protocol.SendReq{}

	rsp := &protocol.AuthRsp{
		Code:    0,
		Message: "成功",
	}

	defer func() {
		b, err := proto.Marshal(rsp)
		if err != nil {
			return
		}

		p.Cmd = uint32(protocol.CmdId_Cmd_SendRsp)
		p.BodyLen = uint32(len(b))
		p.Body = b
		buf := &bytes.Buffer{}
		p.Write(buf)
		c.AsyncWrite(buf.Bytes())
	}()

	if err = proto.Unmarshal(p.Body, req); err != nil {
		rsp.Code = -1
		rsp.Message = "协议解析错误"
		log.Error(err)
		return
	}

	s := h.app.GetSessionByConn(c)
	logicClient := h.app.GetLogicClient()
	r := logic.SendReq{
		ConvType:   req.ConvType,
		MsgType:    req.MsgType,
		From:       req.From,
		To:         req.To,
		Content:    req.Content,
		Extra:      req.Extra,
		ClientTime: req.ClientTime,
		ConnId:     s.ConnId,
	}
	rspL, err := logicClient.SendMsg(context.Background(), &r)
	if err != nil {
		rsp.Code = -1
		rsp.Message = "系统错误"
		log.Error(err)
		return
	}

	rsp.Code = rspL.Code
	rsp.Message = rspL.Message
	return
}
