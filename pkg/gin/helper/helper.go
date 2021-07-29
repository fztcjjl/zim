package helper

import (
	"fmt"
	"github.com/fztcjjl/zim/pkg/gin/context"
	"github.com/gin-gonic/gin"
	"path"
	"regexp"
	"strings"
	"sync"
)

var (
	routerTitle  = &sync.Map{}
	routerRegexp = regexp.MustCompile(`(.*):[^/]+(.*)`)
)

// SetRouterTitle 设定路由标题
func SetRouterTitle(method, router, title string) {
	routerTitle.Store(fmt.Sprintf("%s-%s", method, router), title)
}

// GetRouterTitleAndKey 获取路由标题和键
func GetRouterTitleAndKey(method, router string) (string, string) {
	key := fmt.Sprintf("%s-%s", method, router)
	vv, ok := routerTitle.Load(key)
	if ok {
		return vv.(string), key
	}

	var title string
	routerTitle.Range(func(vk, vv interface{}) bool {
		vkey := vk.(string)
		if !strings.Contains(vkey, "/:") {
			return true
		}

		rkey := "^" + routerRegexp.ReplaceAllString(vkey, "$1[^/]+$2") + "$"
		b, _ := regexp.MatchString(rkey, key)
		if b {
			title = vv.(string)
			key = vkey
		}
		return !b
	})

	return title, key
}

type HandlerFunc func(*context.Context)

func Handle(g *gin.RouterGroup, httpMethod string, relativePath string, handler HandlerFunc, title string) {
	SetRouterTitle(httpMethod, path.Join(g.BasePath(), relativePath), title)
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

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
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
