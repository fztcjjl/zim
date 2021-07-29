package controller

import (
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/logic/handler/http/api"
	"github.com/fztcjjl/zim/pkg/errors"
	"github.com/fztcjjl/zim/pkg/gin/context"
)

type MsgController struct {
}

func (c *MsgController) Send(ctx *context.Context) {
	req := api.SendRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		e := errors.ErrInvalidParam
		e.Detail = err.Error()
		ctx.ResponseError(e)
		return
	}
	log.Debug(req)
}

func (c *MsgController) SyncMsg(ctx *context.Context) {

}
