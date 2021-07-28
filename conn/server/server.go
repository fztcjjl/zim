package server

import (
	"context"
	"fmt"
	"github.com/fztcjjl/tiger/trpc/client"
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/pkg/errors"
	"github.com/fztcjjl/zim/pkg/ztimer"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/panjf2000/gnet/pool/goroutine"
	"google.golang.org/grpc"
	"time"
)

const (
	WsUpgrading = 0
	AuthPending = 1
	Authed      = 2
)

type Server struct {
	opts      Options
	nc        *nats.Conn
	serverId  string
	tcpServer *TcpServer
	wsServer  *WsServer
	timer     *ztimer.Timer
	// TODO: 分桶
	clientManager *ClientManager
	logicClient   logic.LogicClient
	workerPool    *goroutine.Pool
}

func NewServer(opt ...Option) *Server {
	s := new(Server)
	s.opts = NewOptions(opt...)
	s.clientManager = NewClientManager()
	s.workerPool = goroutine.Default()

	if s.opts.TcpPort != "" {
		s.tcpServer = NewTcpServer(s, ":"+s.opts.TcpPort)
	}

	if s.opts.WsPort != "" {
		s.wsServer = NewWsServer(s, ":"+s.opts.WsPort)
	}

	s.timer = ztimer.NewTimer(100, 20)

	cli := client.NewClient(
		"srv.logic",
		client.Registry(s.opts.Registry),
		client.GrpcDialOption(grpc.WithInsecure()),
	)
	if cli == nil {
		log.Fatal("NewClient error")
	}

	s.logicClient = logic.NewLogicClient(cli.GetConn())

	var err error
	if s.nc, err = nats.Connect(s.opts.Nats); err != nil {
		log.Fatal(err)
	}

	return s
}

func (s *Server) GetClientManager() *ClientManager {
	return s.clientManager
}

func (s *Server) GetLogicClient() logic.LogicClient {
	return s.logicClient
}

func (s *Server) GetServerId() string {
	return s.serverId
}

func (s *Server) GetTcpServer() *TcpServer {
	return s.tcpServer
}

func (s *Server) GetWsServer() *WsServer {
	return s.wsServer
}

func (s *Server) GetTimer() *ztimer.Timer {
	return s.timer
}

func (s *Server) Start() {
	go func() {
		s.consume()
	}()
	go func() {
		s.timer.Start()
	}()
	go func() {
		if s.tcpServer != nil {
			if err := s.tcpServer.Start(); err != nil {
				log.Fatal(err)
			}
		}
	}()
	go func() {
		if s.wsServer != nil {
			if err := s.wsServer.Start(); err != nil {
				log.Fatal(err)
			}
		}
	}()
}

func (s *Server) consume() {
	// process push message
	pushMsg := new(logic.PushMsg)
	topic := fmt.Sprintf("zim-push-topic-%s", s.serverId)
	if _, err := s.nc.Subscribe(topic, func(msg *nats.Msg) {

		if err := proto.Unmarshal(msg.Data, pushMsg); err != nil {
			log.Errorf("proto.Unmarshal(%v) error(%v)", msg, err)
			return
		}

		log.Debug("recv a msg", pushMsg)
		for _, v := range pushMsg.ConnIds {
			if client := s.GetClientManager().Get(v); client != nil {
				if client.Conn != nil {
					p := Packet{
						HeaderLen:     20,
						ClientVersion: 1,
						Cmd:           uint32(protocol.CmdId_Cmd_PushMsg),
						Seq:           0,
						BodyLen:       uint32(len(pushMsg.Msg)),
						Body:          pushMsg.Msg,
					}
					client.WritePacket(&p)
				}

			}

		}

	}); err != nil {
		return
	}
}

func (s *Server) OnOpen(client *Client) {
	// 10秒钟之内没有认证成功，关闭连接
	client.TimerTask = s.GetTimer().AfterFunc(time.Second*10, func() {
		client.Close()
	})
}

func (s *Server) OnClose(client *Client) {
	log.Debug("TCP OnClose ...")

	if client.ConnId != "" {
		s.GetClientManager().Remove(client.ConnId)
	}

	if client.Status != Authed {
		return
	}

	s.workerPool.Submit(func() {
		if client != nil {
			logicClient := s.GetLogicClient()
			req := logic.DisconnectReq{
				Uin:      client.Uin,
				Platform: client.Platform,
			}
			logicClient.Disconnect(context.Background(), &req)
		}
	})
}

func (s *Server) OnMessage(data []byte, client *Client) {
	s.workerPool.Submit(func() {
		p := &Packet{}
		if err := p.Read(data); err != nil {
			client.Close()
			return
		}

		if client.Status == AuthPending {
			if err := s.handleAuth(client, p); err != nil {
				client.Close()
			} else {
				client.Status = Authed
			}
		} else {
			s.handleProto(client, p)
		}
	})

}

func (s *Server) handleAuth(client *Client, p *Packet) (err error) {
	log.Info("handleAuth ...")

	req := &protocol.AuthReq{}

	rsp := &protocol.AuthRsp{
		Code:    0,
		Message: "成功",
	}

	defer func() {
		b, err := proto.Marshal(rsp)
		if err != nil {
			log.Error("系统错误")
			return
		}

		p.Cmd = uint32(protocol.CmdId_Cmd_AuthRsp)
		p.BodyLen = uint32(len(b))
		p.Body = b
		client.WritePacket(p)
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

	logicClient := s.GetLogicClient()
	reqL := logic.ConnectReq{
		Uin:      req.Uin,
		Platform: req.Platform,
		Server:   s.GetServerId(),
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
		oldClient := s.GetClientManager().Remove(rspL.KickedConnId)
		if oldClient != nil && oldClient.Conn != nil {
			kick := &protocol.Kick{KickReason: rspL.KickedReason}
			if b, err := proto.Marshal(kick); err != nil {
				oldClient.Write(b)
			}
			log.Info("close old")
			oldClient.Close()
		}
	}

	client.ConnId = rspL.ConnId
	client.Uin = reqL.Uin
	client.Platform = req.Platform
	client.Server = s.GetServerId()
	s.GetClientManager().Add(client)

	// 取消定时任务
	client.TimerTask.Cancel()
	client.TimerTask = nil

	return
}

func (s *Server) handleProto(client *Client, p *Packet) (err error) {
	if p.Cmd == uint32(protocol.CmdId_Cmd_Noop) {
		err = s.handleNoop(client, p)
	} else if p.Cmd == uint32(protocol.CmdId_Cmd_SendReq) {
		err = s.handleSend(client, p)
	} else if p.Cmd == uint32(protocol.CmdId_Cmd_SyncMsgReq) {
		err = s.handleSyncMsg(client, p)
	} else if p.Cmd == uint32(protocol.CmdId_Cmd_MsgAckReq) {
		err = s.handleMsgAckReq(client, p)
	}

	return
}

func (s *Server) handleMsgAckReq(client *Client, p *Packet) (err error) {
	return
}

func (s *Server) handleSyncMsg(client *Client, p *Packet) (err error) {
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
		client.WritePacket(p)
	}()

	if err = proto.Unmarshal(p.Body, req); err != nil {
		return
	}

	logicClient := s.GetLogicClient()
	reqL := logic.SyncMsgReq{
		Uin:    client.Uin,
		ConnId: client.ConnId,
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

func (s *Server) handleNoop(client *Client, p *Packet) (err error) {
	client.WritePacket(p)
	logicClient := s.GetLogicClient()
	req := logic.HeartbeatReq{
		Uin:    client.Uin,
		ConnId: client.ConnId,
		Server: client.Server,
	}
	logicClient.Heartbeat(context.Background(), &req)

	return
}

func (s *Server) handleSend(client *Client, p *Packet) (err error) {
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
		client.WritePacket(p)
	}()

	if err = proto.Unmarshal(p.Body, req); err != nil {
		rsp.Code = -1
		rsp.Message = "协议解析错误"
		log.Error(err)
		return
	}

	logicClient := s.GetLogicClient()
	r := logic.SendReq{
		ConvType:   req.ConvType,
		MsgType:    req.MsgType,
		From:       req.From,
		To:         req.To,
		Content:    req.Content,
		Extra:      req.Extra,
		ClientTime: req.ClientTime,
		ConnId:     client.ConnId,
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
