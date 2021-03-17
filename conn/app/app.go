package app

import (
	"bytes"
	"fmt"
	"github.com/fztcjjl/zim/pkg/ztimer"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fztcjjl/tiger/trpc/client"
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/tiger/trpc/registry"
	"github.com/fztcjjl/tiger/trpc/registry/etcd"
	"github.com/fztcjjl/zim/api/logic"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/panjf2000/gnet"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	WsUpgrading = 0
	AuthPending = 1
	Authed      = 2
)

type App struct {
	sync.Mutex
	opts      Options
	config    *Config
	nc        *nats.Conn
	serverId  string
	tcpServer *TcpServer
	wsServer  *WsServer
	timer     *ztimer.Timer
	// TODO: 分桶
	idSessions  map[string]*Session
	logicClient logic.LogicClient
}

type Session struct {
	Status    int
	TimerTask *ztimer.TimerTask

	ConnId   string
	Conn     gnet.Conn
	Uin      string
	Platform string
	Server   string
}

func NewApp(opt ...Option) *App {
	app := new(App)
	app.idSessions = make(map[string]*Session)
	app.loadConfig()
	app.initLogger()
	//app.initTracer()
	options := newOptions(opt...)
	app.opts = options

	app.serverId = app.config.GetString("app.server_id")

	if app.config.GetString("tcp.addr") != "" {
		app.tcpServer = NewTcpServer(app, app.config.GetString("tcp.addr"))
	}

	if app.config.GetString("ws.addr") != "" {
		app.wsServer = NewWsServer(app, app.config.GetString("ws.addr"))
	}

	app.timer = ztimer.NewTimer(100, 20)
	return app
}

func (a *App) AddSession(s *Session) {
	a.Lock()
	a.idSessions[s.ConnId] = s
	a.Unlock()
}

func (a *App) DelSessionByConnId(id string) (s *Session) {
	a.Lock()
	s = a.idSessions[id]
	if s != nil {
		delete(a.idSessions, id)
	}

	a.Unlock()

	return
}

func (a *App) GetSessionByConnId(id string) (s *Session) {
	a.Lock()
	s = a.idSessions[id]
	a.Unlock()
	return
}

func (a *App) GetLogicClient() logic.LogicClient {
	return a.logicClient
}

func (a *App) GetServerId() string {
	return a.serverId
}

func (a *App) GetTcpServer() *TcpServer {
	return a.tcpServer
}

func (a *App) GetWsServer() *WsServer {
	return a.wsServer
}

func (a *App) Name() string {
	return a.config.GetString("app.name")
}

func (a *App) GetTimer() *ztimer.Timer {
	return a.timer
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
	var r registry.Registry
	addrs := a.config.GetStringSlice("etcd")
	if len(addrs) > 0 {
		r = etcd.NewRegistry(registry.Addrs(addrs...))
	}

	cli := client.NewClient(
		"srv.logic",
		client.Registry(r),
		client.GrpcDialOption(grpc.WithInsecure()),
	)
	if cli == nil {
		log.Fatal("NewClient error")
	}

	a.logicClient = logic.NewLogicClient(cli.GetConn())
	var err error
	addr := a.config.GetString("nats.addr")
	if a.nc, err = nats.Connect(addr); err != nil {
		log.Fatal(err)
	}

	go func() {
		a.consume()
	}()

	go func() {
		if a.tcpServer != nil {
			if err := a.tcpServer.Start(); err != nil {
				log.Fatal(err)
			}
		}
	}()
	go func() {
		if a.wsServer != nil {
			if err := a.wsServer.Start(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	select {
	case <-ch:
	}

}

func (a *App) consume() {
	// process push message
	pushMsg := new(logic.PushMsg)
	topic := fmt.Sprintf("zim-push-topic-%s", a.serverId)
	if _, err := a.nc.Subscribe(topic, func(msg *nats.Msg) {

		if err := proto.Unmarshal(msg.Data, pushMsg); err != nil {
			log.Errorf("proto.Unmarshal(%v) error(%v)", msg, err)
			return
		}

		log.Debug("recv a msg", pushMsg)
		for _, v := range pushMsg.ConnIds {
			if s := a.idSessions[v]; s != nil {
				if s.Conn != nil {
					p := protocol.Proto{
						HeaderLen:     20,
						ClientVersion: 1,
						Cmd:           uint32(protocol.CmdId_Cmd_PushMsg),
						Seq:           0,
						BodyLen:       uint32(len(pushMsg.Msg)),
						Body:          pushMsg.Msg,
					}
					buf := &bytes.Buffer{}
					p.Write(buf)
					s.Conn.AsyncWrite(buf.Bytes())
				}

			}

		}

	}); err != nil {
		return
	}
}
