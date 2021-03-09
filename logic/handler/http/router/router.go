package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Setup(engine *gin.Engine) {
	//initSentinel()
	engine.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})

	engine.Use(gin.Recovery())
	engine.Use(requestid.New())

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST, OPTIONS, GET, PUT, PATCH, DELETE"},
		AllowHeaders: []string{"*"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	//engine.Use(
	//	middleware.Sentinel(
	//		middleware.WithResourceExtractor(func(ctx *gin.Context) string {
	//			return ctx.GetHeader("X-Real-IP")
	//		}),
	//		middleware.WithBlockFallback(func(ctx *gin.Context) {
	//			ctx.AbortWithStatusJSON(400, map[string]interface{}{
	//				"code":    9999,
	//				"message": "服务器忙，请稍后重试",
	//			})
	//		}),
	//	),
	//)
	//engine.Use(middleware.CORSMiddleware())

	//apiPrefixes := []string{"/mobile/", "/swagger/"}
	//engine.Use(middleware.StaticFile("www", apiPrefixes...))
	// 注册/mobile/v1路由
	//api.RegisterV1(engine)
	//RegisterUmsgRouter(engine)
	//RegisterNotifyRouter(engine)
}
