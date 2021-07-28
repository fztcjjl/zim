package rpc

import (
	"context"
	"fmt"
	"time"

	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/logic/app"
	"github.com/fztcjjl/zim/logic/dao"
	"github.com/fztcjjl/zim/logic/model"
	"github.com/fztcjjl/zim/pkg/util"
	"github.com/google/uuid"
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
	if req.ConvType == 1 {
		rsp, err = l.sendC2C(ctx, req)
	} else if req.ConvType == 2 {
		rsp, err = l.sendC2G(ctx, req)
	}
	return
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
