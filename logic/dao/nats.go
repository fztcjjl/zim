package dao

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/fztcjjl/zim/api/logic"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

var (
	nc     *Nats
	onceNc sync.Once
)

type Nats struct {
	nc *nats.Conn
}

func getNats() *Nats {
	onceNc.Do(func() {
		nc = new(Nats)
	})
	return nc
}

func SetupNats(addr string) {
	r := getNats()
	nc, err := nats.Connect(addr)
	if err != nil {
		log.Panic(err)
	}
	r.nc = nc
	//r
	//r.client = redis.NewClient(&redis.Options{
	//	Addr:     addr,
	//	Password: password,
	//	DB:       0,
	//})
}

func GetNatsConn() *nats.Conn {
	return getNats().nc
	//onceNc.Do(func() {
	//	nc, _ = nats.Connect(nats.DefaultURL)
	//})
	//
	//return nc
}

// PushMsg push a message to databus.
func PushMsg(ctx context.Context, server string, connIds []string, msg []byte) (err error) {
	pushMsg := &logic.PushMsg{
		//Cmd:     int32(protocol.CmdId_Cmd_PushMsg),
		Server:  server,
		ConnIds: connIds,
		Msg:     msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return
	}

	topic := fmt.Sprintf("zim-push-topic-%s", server)

	fmt.Println(topic)
	m := &nats.Msg{
		Subject: topic,
		Reply:   "ack",
		Data:    b,
		Sub:     nil,
	}

	GetNatsConn().PublishMsg(m)
	return
}
