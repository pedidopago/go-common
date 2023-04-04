package redis

import (
	"context"
	"time"
)

type Client interface {
	Get(ctx context.Context, key string) StringCmder
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) StatusCmder
}

type StatusCmder interface {
	Name() string
	FullName() string
	Args() []interface{}
	Err() error
	Val() string
	Result() (string, error)
	String() string
}

type StringCmder interface {
	Name() string
	FullName() string
	Args() []interface{}
	Err() error
	Val() string
	Result() (string, error)
	Bytes() ([]byte, error)
	Int() (int, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Time() (time.Time, error)
	Scan(val interface{}) error
	String() string
}

// type IntCmder interface {
// 	SetVal(val int64)
// 	Val() int64
// 	Result() (int64, error)
// 	Uint64() (uint64, error)
// 	String() string
// }

// type PubSuber interface {
// 	String() string
// }

// func (*PubSub).Channel(opts ...ChannelOption) <-chan *Message
// func (*PubSub).ChannelSize(size int) <-chan *Message
// func (*PubSub).ChannelWithSubscriptions(_ context.Context, size int) <-chan interface{}
// func (*PubSub).Close() error
// func (*PubSub).PSubscribe(ctx context.Context, patterns ...string) error
// func (*PubSub).PUnsubscribe(ctx context.Context, patterns ...string) error
// func (*PubSub).Ping(ctx context.Context, payload ...string) error
// func (*PubSub).Receive(ctx context.Context) (interface{}, error)
// func (*PubSub).ReceiveMessage(ctx context.Context) (*Message, error)
// func (*PubSub).ReceiveTimeout(ctx context.Context, timeout time.Duration) (interface{}, error)
// func (*PubSub).String() string
// func (*PubSub).Subscribe(ctx context.Context, channels ...string) error
// func (*PubSub).Unsubscribe(ctx context.Context, channels ...string) error
