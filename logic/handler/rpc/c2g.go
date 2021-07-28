package rpc

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"time"

	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/fztcjjl/zim/logic/dao"
	"github.com/fztcjjl/zim/logic/model"
	"github.com/fztcjjl/zim/pkg/idgen"
	"github.com/golang/protobuf/proto"
)

func (l *Logic) sendC2G(ctx context.Context, req *logic.SendReq) (rsp *logic.SendRsp, err error) {
	log.Debug(req)

	db := dao.GetDB()
	now := time.Now()

	msg := model.ImMsgSend{
		MsgId:      idgen.Next(),
		ConvType:   int(req.ConvType),
		Content:    req.Content,
		Extra:      req.Extra,
		Type:       int(req.MsgType),
		Sender:     req.Sender,
		Target:     req.Target,
		AtUserList: "",
	}

	b, _ := json.Marshal(req.AtUserList)
	msg.AtUserList = string(b)

	if err = db.Create(&msg).Error; err != nil {
		log.Error(err)
		return
	}

	var members []*model.ImGroupMember
	cond := model.ImGroupMember{
		GroupId: req.Target,
	}
	if e := db.Where(&cond).Find(&members).Error; e != nil {
		log.Error(e)
		err = e
		return
	}

	msgr := model.ImMsgRecv{
		MsgId:      msg.MsgId,
		ConvType:   msg.ConvType,
		Content:    msg.Content,
		Extra:      msg.Extra,
		Type:       msg.Type,
		Sender:     msg.Sender,
		Target:     msg.Target,
		AtUserList: msg.AtUserList,
	}

	for _, m := range members {
		msgr.Receiver = m.Member
		if err := db.Create(&msg).Error; err != nil {
			log.Error(err)
			continue
		}
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
		ConvType:   req.ConvType,
		Type:       req.MsgType,
		Content:    req.Content,
		Sender:     req.Sender,
		Target:     req.Target,
		Extra:      req.Extra,
		SendTime:   now.Unix(),
		AtUserList: req.AtUserList,
	}
	b, e := proto.Marshal(&p)
	if e != nil {
		return
	}

	for _, m := range members {
		if m.Id == cast.ToInt64(req.Sender) {
			//var excludeConnId string
			//if conn := dao.GetConnByPlatform(req.Target, plat); conn != nil {
			//	excludeConnId = conn.ConnId
			//	um.PushByUin(context.Background(), m.Member, excludeConnId, b)
			//}
			pushByUin(ctx, m.Member, "", b)
		} else {
			pushByUin(ctx, m.Member, "", b)
		}

	}

	return
}
