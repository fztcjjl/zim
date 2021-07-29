package context

import (
	"github.com/fztcjjl/zim/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
	"time"
)

type Context struct {
	*gin.Context
}

func New(ctx *gin.Context) *Context {
	return &Context{Context: ctx}
}

func (ctx *Context) Response(obj interface{}) {
	m := make(map[string]interface{})
	m["code"] = 0
	m["data"] = obj
	m["ts"] = time.Now().UnixNano()
	ctx.response(http.StatusOK, m)
}

func (ctx *Context) ResponseOK() {
	ctx.Response(nil)
}

func (ctx *Context) ResponseError(err error) {
	ce := errors.Parse(err.Error())

	m := make(map[string]interface{})
	m["code"] = ce.Code
	if ce.Message != "" {
		m["message"] = ce.Message
	}
	if ce.Detail != "" {
		m["detail"] = ce.Detail
	}
	m["ts"] = time.Now().UnixNano()
	if ce.Code == -1 {
		ctx.response(500, m)
		return
	}
	ctx.response(499, m)

	return
}

func (ctx *Context) response(status int, obj interface{}) {
	ctx.JSON(status, obj)
	ctx.Abort()
}

var (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

//func (c *Context) GetAppHeader() *typ.AppHeader {
//	//log.Debug(c.gctx.Request.Header)
//	h := &typ.AppHeader{}
//	c.gctx.ShouldBindHeader(h)
//	return h
//}
//
//func (c *Context) GetWsHeader() *typ.WsHeader {
//	h := &typ.WsHeader{}
//	c.gctx.ShouldBindHeader(h)
//	return h
//}

func (ctx *Context) GetForm() url.Values {
	ctx.Request.ParseForm()
	return ctx.Request.PostForm
}

func (ctx *Context) GetUidStr() string {
	return cast.ToString(ctx.GetUid())
}

func (ctx *Context) GetUid() int64 {
	//if v, exists := ctx.Get(constant.LuoboAccountKey); !exists {
	//	return 0
	//} else {
	//	if acc, ok := v.(jwt.Account); !ok {
	//		return 0
	//	} else {
	//		return acc.Uid
	//	}
	//}
	return 0
}

func (ctx *Context) GetPlatform() int {
	return 0
	//if v, exists := c.gctx.Get(constant.LuoboAccountKey); !exists {
	//	return 0
	//} else {
	//	if acc, ok := v.(jwt.Account); !ok {
	//		return 0
	//	} else {
	//		return acc.Platform
	//	}
	//}
}

func (ctx *Context) GetDeviceName() string {
	//if v, exists := c.gctx.Get(constant.LuoboAccountKey); !exists {
	//	return ""
	//} else {
	//	if acc, ok := v.(jwt.Account); !ok {
	//		return ""
	//	} else {
	//		return acc.DeviceName
	//	}
	//}

	return ""
}

func (ctx *Context) GetToken() string {
	return ctx.GetHeader("Token")
}

func (c *Context) GetPage() int {
	if v := c.Query("page"); len(v) > 0 {
		if n := cast.ToInt(v); n > 0 {
			return n
		}
	}

	return 1
}

func (c *Context) GetPageSize() int {
	if v := c.Query("per_page"); len(v) > 0 {
		if n := cast.ToInt(v); n > 0 {
			if n > MaxPageSize {
				return MaxPageSize
			}
			return n
		}
	}

	return DefaultPageSize
}

type list struct {
	List    interface{} `json:"list,omitempty"`
	Total   int         `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	PerPage int         `json:"per_page,omitempty"`
	//Pagination *pagination `json:"pagination,omitempty"`
}

type pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (c *Context) ResponseList(obj interface{}) {
	c.Response(list{List: obj})
}

func (c *Context) ResponsePage(total int, obj interface{}) {
	c.Response(list{
		List:    obj,
		Total:   total,
		Page:    c.GetPage(),
		PerPage: c.GetPageSize(),
		//Pagination: &pagination{
		//	Total:   total,
		//	Page:    c.GetPage(),
		//	PerPage: c.GetPageSize(),
		//},
	})
}
