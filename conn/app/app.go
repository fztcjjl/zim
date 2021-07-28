package app

import (
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/tiger/trpc/registry"
	"github.com/fztcjjl/tiger/trpc/registry/etcd"
	"github.com/fztcjjl/zim/conn/server"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	opts   Options
	config *Config
	server *server.Server
}

func NewApp(opt ...Option) *App {
	app := new(App)
	app.opts = newOptions(opt...)
	app.loadConfig()
	app.initLogger()

	var r registry.Registry
	addrs := app.config.GetStringSlice("etcd")
	if len(addrs) > 0 {
		r = etcd.NewRegistry(registry.Addrs(addrs...))
	}

	app.server = server.NewServer(
		server.Registry(r),
		server.WithServerId(app.config.GetString("server.id")),
		server.WithTcpPort(app.config.GetString("server.tcp_port")),
		server.WithWsPort(app.config.GetString("server.ws_port")),
		server.WithNats(app.config.GetString("nats.addr")),
	)

	return app
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

func (a *App) GetConfig() *Config {
	return a.config
}

func (a *App) Run() {
	a.server.Start()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	select {
	case <-ch:
	}

}
