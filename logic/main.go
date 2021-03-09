package main

import (
	"github.com/fztcjjl/tiger/trpc/web"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/logic/app"
	"github.com/fztcjjl/zim/logic/dao"
	"github.com/fztcjjl/zim/logic/handler/http"
	"github.com/fztcjjl/zim/logic/handler/rpc"
)

func main() {
	a := app.NewApp(app.WithHttp(true))
	srv := a.GetServer()
	webSrv := a.GetWebServer()
	webSrv.Init(web.Handler(http.Handler()))
	logic.RegisterLogicServer(srv.Server(), &rpc.Logic{App: a})

	setup(a)
	a.Run()
}

func setup(a *app.App) {
	c := dao.Config{
		DriverName:     "mysql",
		DataSourceName: a.GetConfig().GetString("mysql.data_source"),
		MaxIdleConn:    a.GetConfig().GetInt("mysql.max_idle"),
		MaxOpenConn:    a.GetConfig().GetInt("mysql.max_conn"),
	}
	dao.Setup(c)
	dao.SetupRedis(
		a.GetConfig().GetString("redis.addr"),
		a.GetConfig().GetString("redis.password"),
		a.GetConfig().GetInt("redis.db"),
	)
	dao.SetupNats(a.GetConfig().GetString("nats.addr"))
}
