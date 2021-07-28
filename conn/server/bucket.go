package server

import (
	"sync"
)

type Bucket struct {
	sync.RWMutex
	clients map[string]*Client
}
