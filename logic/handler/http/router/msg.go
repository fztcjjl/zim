package router

import (
	"github.com/fztcjjl/zim/logic/handler/http/controller"
	"github.com/fztcjjl/zim/pkg/gin/helper"
	"github.com/gin-gonic/gin"
)

func RegisterMsgRouter(g *gin.RouterGroup, c *controller.MsgController) {
	//helper.POST(g, "/im/login", c.Login, "登录")
	helper.POST(g, "/im/send", c.Send, "发送消息")
	helper.POST(g, "/im/sync_msg", c.SyncMsg, "同步消息")
}
