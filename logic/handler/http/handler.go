package http

import (
	"github.com/fztcjjl/zim/logic/handler/http/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler() http.Handler {
	engine := gin.New()
	router.Setup(engine)
	return engine
}
