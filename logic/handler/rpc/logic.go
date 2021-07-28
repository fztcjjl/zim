package rpc

import (
	"context"
	"fmt"
	"time"

	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/logic/app"
	"github.com/fztcjjl/zim/logic/dao"
	"github.com/fztcjjl/zim/logic/model"
	"github.com/fztcjjl/zim/pkg/idgen"
	"github.com/fztcjjl/zim/pkg/util"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Logic struct {
	App *app.App
}

func (l *Logic) Connect(ctx context.Context, req *logic.ConnectReq) (rsp *logic.ConnectRsp, err error) {
	rsp = &logic.ConnectRsp{}
	if req.Platform == "android" || req.Platform == "iOS" {

	}
	if info := dao.GetConnByPlatform(ctx, req.Uin, req.Platform); info != nil {
		rsp.KickedConnId = info.ConnId
		rsp.KickedReason = fmt.Sprintf("您的账号已在设备%s上登录", req.Device)
	}
	connId := uuid.New().String()

	info := &dao.ConnInfo{
		ConnId:         connId,
		Platform:       req.Platform,
		LoginTime:      time.Now().Unix(),
		DisconnectTime: 0,
		Device:         req.Device,
		Status:         dao.Online,
		Server:         req.Server,
	}
	dao.AddConn(ctx, req.Uin, info)
	rsp.ConnId = connId

	// TODO: token验证

	return
}

func (l *Logic) Disconnect(ctx context.Context, req *logic.DisconnectReq) (rsp *logic.DisconnectRsp, err error) {
	info := dao.GetConnByPlatform(ctx, req.Uin, req.Platform)
	if info == nil {
		return
	}

	info.DisconnectTime = time.Now().Unix()
	info.Status = dao.PushOnline

	// 覆盖
	if err = dao.AddConn(ctx, req.Uin, info); err != nil {
		return
	}
	rsp = &logic.DisconnectRsp{}
	return
}

func (l *Logic) Heartbeat(ctx context.Context, req *logic.HeartbeatReq) (rsp *logic.HeartbeatRsp, err error) {
	if err = dao.ExpireConn(ctx, req.Uin); err != nil {
		return
	}

	rsp = &logic.HeartbeatRsp{}
	return
}

func (l *Logic) SendMsg(ctx context.Context, req *logic.SendReq) (rsp *logic.SendRsp, err error) {
	return l.sendC2C(ctx, req)
}

func (l *Logic) SyncMsg(ctx context.Context, req *logic.SyncMsgReq) (rsp *logic.SyncMsgRsp, err error) {
	log.Debug("offset=%d", req.Offset)
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	var rows []*model.ImMsgRecv
	db := dao.GetDB()
	var where string
	if req.Offset == 0 {
		where = "receiver=? AND msg_id>? AND DATE_SUB(CURDATE(), INTERVAL 7 DAY) <= created_at"
	} else {
		where = "receiver=? AND msg_id>?"
	}
	if err = db.Table("im_msg_recv").Where(where, req.Uin, req.Offset).
		Order("msg_id ASC").Limit(int(req.Limit)).Find(&rows).Error; err != nil {
		return
	}

	rsp = &logic.SyncMsgRsp{}
	for _, v := range rows {
		msg := logic.Msg{
			Id:       v.MsgId,
			ConvType: int32(v.ConvType),
			Type:     int32(v.Type),
			Content:  v.Content,
			Sender:   v.Sender,
			Target:   v.Target,
			Extra:    v.Extra,
			SendTime: util.TimeToUnix(v.CreatedAt),
		}
		rsp.List = append(rsp.List, &msg)
	}

	return
}

func (l *Logic) MsgAck(ctx context.Context, req *logic.MsgAckReq) (rsp *logic.MsgAckRsp, err error) {
	return
}

func (l *Logic) sendC2C(ctx context.Context, req *logic.SendReq) (rsp *logic.SendRsp, err error) {
	log.Debug(req)

	db := dao.GetDB()
	msg := model.ImMsgSend{}

	now := time.Now()
	err = db.Transaction(func(tx *gorm.DB) error {

		msg = model.ImMsgSend{
			MsgId:      idgen.Next(),
			ConvType:   int(req.ConvType),
			Content:    req.Content,
			Extra:      req.Extra,
			Type:       int(req.MsgType),
			Sender:     req.Sender,
			Target:     req.Target,
			AtUserList: "",
		}

		if err := tx.Create(&msg).Error; err != nil {
			log.Error(err)
			return err
		}

		msgr := model.ImMsgRecv{
			MsgId:      msg.MsgId,
			ConvType:   msg.ConvType,
			Content:    msg.Content,
			Extra:      msg.Extra,
			Type:       msg.Type,
			Sender:     msg.Sender,
			Target:     msg.Target,
			Receiver:   msg.Target,
			AtUserList: "",
		}

		var msgs []model.ImMsgRecv
		// 给自己的收件箱也插入一条消息，为了多端同步
		msgr2 := msgr
		msgr2.Receiver = msg.Sender
		msgs = append(msgs, msgr, msgr2)

		if err := tx.Create(&msgs).Error; err != nil {
			log.Error(err)
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	rsp = &logic.SendRsp{
		Code:     0,
		Message:  "",
		Id:       msg.MsgId,
		SendTime: now.Unix(),
		Seq:      0,
	}

	p := protocol.Msg{
		Id:         msg.MsgId,
		ConvType:   int32(req.ConvType),
		Type:       int32(req.MsgType),
		Content:    req.Content,
		Sender:     req.Sender,
		Target:     req.Target,
		Extra:      req.Extra,
		SendTime:   now.Unix(),
		AtUserList: nil,
	}

	b, _ := proto.Marshal(&p)

	pushByUin(ctx, req.Sender, "", b)
	pushByUin(ctx, req.Target, "", b)

	return
}

func pushByUin(ctx context.Context, uin string, exclude string, msg []byte) {
	log.Info(uin)
	receivers, err := dao.GetConnByUin(ctx, uin)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(receivers)
	for server, conns := range receivers {
		var connIds []string
		for _, conn := range conns {
			if exclude != "" && conn.ConnId == exclude {
				continue
			}
			if conn.Status == dao.Online {
				connIds = append(connIds, conn.ConnId)
			} else if conn.Status == dao.PushOnline {
				if time.Since(time.Unix(conn.DisconnectTime, 0)) < time.Duration(dao.PushOnlineKeepDays*24)*time.Hour {
					// TODO: 推送通道
					// TODO: 判断手机类型,走不同推送通道
				}
			}
		}

		log.Info(connIds)

		dao.PushMsg(ctx, server, connIds, msg)
	}
}
