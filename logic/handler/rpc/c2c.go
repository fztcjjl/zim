package rpc

import (
	"context"
	"time"

	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/logic/dao"
	"github.com/fztcjjl/zim/logic/model"
	"github.com/fztcjjl/zim/pkg/idgen"
	"github.com/golang/protobuf/proto"
	"gorm.io/gorm"
)

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
