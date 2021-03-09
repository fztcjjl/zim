package app

import (
	"github.com/panjf2000/gnet"
	"sync"
)

type Bucket struct {
	conns map[string]*gnet.Conn
	sync.RWMutex
}
