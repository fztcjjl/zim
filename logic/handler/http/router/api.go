package router

import (
	"github.com/fztcjjl/zim/logic/handler/http/controller"
	"github.com/fztcjjl/zim/logic/handler/http/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	g := r.Group(
		"/api/v1/",
		middleware.CheckSign(),
	)

	RegisterMsgRouter(g, &controller.MsgController{})
	RegisterConversationRouter(g, &controller.ConversationController{})
	RegisterGroupRouter(g, &controller.GroupController{})
	RegisterProfileRouter(g, &controller.ProfileController{})
	RegisterSnsRouter(g, &controller.SnsController{})
}
