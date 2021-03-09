package router

import (
	"github.com/fztcjjl/zim/logic/handler/http/controller"
	"github.com/gin-gonic/gin"
)

func RegisterProfileRouter(g *gin.RouterGroup, c *controller.ProfileController) {
	//helper.POST(g, "/profile/portrait_get", c.PortraitGet, "拉取资料")
	//helper.POST(g, "/profile/portrait_set", c.PortraitSet, "设置资料")
}
