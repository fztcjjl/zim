package app

import (
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/tiger/trpc/registry"
	"github.com/fztcjjl/tiger/trpc/registry/etcd"
	"github.com/fztcjjl/tiger/trpc/server"
	"github.com/fztcjjl/tiger/trpc/web"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	opts      Options
	config    *Config
	server    *server.Server
	webServer *web.Server
}

func NewApp(opt ...Option) *App {
	app := new(App)
	app.loadConfig()
	app.initLogger()
	//app.initTracer()
	options := newOptions(opt...)
	app.opts = options

	var r registry.Registry
	addrs := app.config.GetStringSlice("etcd")
	if len(addrs) > 0 {
		r = etcd.NewRegistry(registry.Addrs(addrs...))
	}
	name := app.config.GetString("app.name")
	version := app.config.GetString("app.version")
	if app.opts.EnableHttp {
		app.webServer = web.NewServer(
			web.Name("web."+name),
			web.Version(version),
			web.Registry(r),
		)
	}

	rpcAddr := app.config.GetString("rpc.addr")
	app.server = server.NewServer(
		server.Name("srv."+name),
		server.Version(version),
		server.Registry(r),
		server.Address(rpcAddr),
		server.Interceptors(
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			//grpc_zap.UnaryServerInterceptor(zap.Logger()),
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	return app
}

func (a *App) GetServer() *server.Server {
	return a.server
}

func (a *App) GetWebServer() *web.Server {
	return a.webServer
}

func (a *App) Name() string {
	return a.config.GetString("app.name")
}

func (a *App) Init(opt ...Option) {
	for _, o := range opt {
		o(&a.opts)
	}
}

func (a *App) loadConfig() {
	v := viper.New()
	v.AddConfigPath("conf")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
		return
	}

	a.config = &Config{Viper: v}
	return
}

func (a *App) initLogger() {
	lvl, _ := log.GetLevel(a.config.Viper.GetString("log.level"))
	log.Init(log.WithLevel(lvl))
}

//func (a *App) initTracer() {
//	n := a.config.GetString("app.name")
//	addr := a.config.GetString("jaeger.addr")
//	trace.Init(n, addr)
//}

func (a *App) GetConfig() *Config {
	return a.config
}

func (a *App) Run() {
	log.Infof("Starting [service] %s", a.Name())
	if a.server != nil {
		if err := a.server.Start(); err != nil {
			log.Fatal(err)
		}
	}

	if a.webServer != nil {
		if err := a.webServer.Start(); err != nil {
			log.Fatal(err)
		}
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-a.opts.Context.Done():
	}

	if a.server != nil {
		a.server.Stop()
	}

	if a.webServer != nil {
		a.webServer.Stop()
	}

}
