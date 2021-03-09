package helper

import (
	"github.com/fztcjjl/zim/pkg/context"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*context.Context)

func Handle(g *gin.RouterGroup, httpMethod string, relativePath string, handler HandlerFunc, title string) {
	//context.SetRouterTitle(httpMethod, path.Join(g.BasePath(), relativePath), title)
	g.Handle(httpMethod, relativePath, func(c *gin.Context) {
		handler(context.New(c))
	})
}

func GET(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "GET", relativePath, handler, title)
}

func POST(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "POST", relativePath, handler, title)
}

func DELETE(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "DELETE", relativePath, handler, title)
}

func PATCH(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "PATCH", relativePath, handler, title)
}

func PUT(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "PUT", relativePath, handler, title)
}

func OPTIONS(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "OPTIONS", relativePath, handler, title)
}

func HEAD(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "HEAD", relativePath, handler, title)
}

func Any(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "GET", relativePath, handler, title)
	Handle(g, "POST", relativePath, handler, title)
	Handle(g, "PUT", relativePath, handler, title)
	Handle(g, "PATCH", relativePath, handler, title)
	Handle(g, "HEAD", relativePath, handler, title)
	Handle(g, "OPTIONS", relativePath, handler, title)
	Handle(g, "DELETE", relativePath, handler, title)
	Handle(g, "CONNECT", relativePath, handler, title)
	Handle(g, "TRACE", relativePath, handler, title)
}
