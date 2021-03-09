package rpc

import (
	"context"
	"database/sql"
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
	"github.com/spf13/cast"
	"github.com/zentures/cityhash"
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

	db := dao.GetDB()
	// 同步新消息
	index := cast.ToUint32(req.Uin)
	if index != 0 {
		index = index % 8
	} else {
		index = cityhash.CityHash32([]byte(req.Uin), uint32(len(req.Uin))) % 8
	}
	sql := "select * from im_msg_recv_%02d where delivered=0 and `to`=? and seq>? order by seq asc limit ?"

	result := db.Raw(fmt.Sprintf(sql, index), req.Uin, req.Offset, req.Limit)
	if err = result.Error; err != nil {
		log.Error(err)
		return
	}

	rows, err := result.Rows()
	if err != nil {
		log.Error(err)
		return
	}

	defer rows.Close()

	rsp = &logic.SyncMsgRsp{}
	for rows.Next() {
		v := model.ImMsgRecv00{}
		if err = db.ScanRows(rows, &v); err != nil {
			return nil, err
		}
		msg := logic.Msg{
			Id:       v.Id,
			ConvType: int32(v.ConvType),
			Type:     int32(v.Type),
			Content:  v.Content,
			From:     v.From,
			To:       v.To,
			Extra:    v.Extra,
			SendTime: util.TimeToUnix(v.CreatedAt),
			Seq:      v.Seq,
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
	fromSeq, err := incr(1, req.From)
	toSeq, err := incr(1, req.To)
	m1 := model.ImMsgRecv00{}

	now := time.Now()
	err = db.Transaction(func(tx *gorm.DB) error {
		sql := "insert into im_msg_recv_%02d (id,conv_type,type,content,extra,created_at,updated_at,`from`,`to`,target,seq,client_time) values(?,?,?,?,?,?,?,?,?,?,?,?)"

		m1.Id = idgen.Next()
		m1.ConvType = int(req.ConvType)
		m1.Type = int(req.MsgType)
		m1.Content = req.Content
		m1.Extra = req.Extra
		m1.From = req.From
		m1.To = req.From
		m1.Target = req.To
		m1.Delivered = 1
		m1.Seq = fromSeq

		index := cast.ToUint32(m1.To)
		if index != 0 {
			index = cast.ToUint32(m1.To) % 8
		} else {
			index = cityhash.CityHash32([]byte(m1.To), uint32(len(m1.To))) % 8
			log.Info(m1.To, index)
		}

		err = tx.Exec(fmt.Sprintf(sql, index),
			m1.Id, m1.ConvType, m1.Type, m1.Content,
			m1.Extra, util.TimeFormat(now), util.TimeFormat(now), m1.From, m1.To, m1.Target, m1.Seq, req.ClientTime).Error
		if err != nil {
			log.Error(err)
			return err
		}

		m1.Id = idgen.Next()
		m1.To = req.To
		m1.Seq = toSeq
		index = cast.ToUint32(m1.To)
		if index != 0 {
			index = cast.ToUint32(m1.To) % 8
		} else {
			index = cityhash.CityHash32([]byte(m1.To), uint32(len(m1.To))) % 8
			log.Info(m1.To, index)
		}
		err = tx.Exec(fmt.Sprintf(sql, index),
			m1.Id, m1.ConvType, m1.Type, m1.Content,
			m1.Extra, util.TimeFormat(now), util.TimeFormat(now), m1.From, m1.To, m1.Target, m1.Seq, req.ClientTime).Error
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	log.Info("LLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN")
	p := protocol.Msg{
		Id:       m1.Id,
		ConvType: req.ConvType,
		Type:     req.MsgType,
		Content:  req.Content,
		From:     req.From,
		To:       req.To,
		Extra:    req.Extra,
		SendTime: now.Unix(),
		Seq:      fromSeq,
	}
	b, _ := proto.Marshal(&p)
	pushByUin(ctx, req.From, "", b)

	b, _ = proto.Marshal(&p)
	p.Seq = toSeq
	pushByUin(ctx, req.To, "", b)

	rsp = &logic.SendRsp{
		Id:       m1.Id,
		Seq:      m1.Seq,
		SendTime: now.Unix(),
	}

	return
}
func incr(objType int, objId string) (seq int64, err error) {
	db := dao.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Raw("select seq from seq where obj_type=? and obj_id=? for update", objType, objId).Row().Scan(&seq)
		if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			return err
		}

		if err == sql.ErrNoRows {
			err = tx.Exec("insert into seq (obj_type,obj_id,seq) values (?,?,?)", objType, objId, seq+1).Error
			if err != nil {
				log.Error(err)
				return err
			}
		} else {
			err = tx.Exec("update seq set seq = seq + 1 where obj_type = ? and obj_id = ?", objType, objId).Error
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	})

	seq = seq + 1
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
