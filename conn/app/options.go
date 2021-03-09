package app

import (
	"context"
)

type Options struct {
	EnableTcp       bool
	EnableWebSocket bool
	Context         context.Context
}

func newOptions(opt ...Option) Options {
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

func WithTcp(enable bool) Option {
	return func(o *Options) {
		o.EnableTcp = enable
	}
}

func WithWebsocket(enable bool) Option {
	return func(o *Options) {
		o.EnableWebSocket = enable
	}
}
