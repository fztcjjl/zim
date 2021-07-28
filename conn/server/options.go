package server

import (
	"context"
	"github.com/fztcjjl/tiger/trpc/registry"
)

type Options struct {
	ServerId string
	TcpPort  string
	WsPort   string
	Context  context.Context
	Registry registry.Registry
	Nats     string
}

func NewOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

type Option func(*Options)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func WithTcpPort(port string) Option {
	return func(o *Options) {
		o.TcpPort = port
	}
}

func WithWsPort(port string) Option {
	return func(o *Options) {
		o.WsPort = port
	}
}

func WithServerId(id string) Option {
	return func(o *Options) {
		o.ServerId = id
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func WithNats(addr string) Option {
	return func(o *Options) {
		o.Nats = addr
	}
}
