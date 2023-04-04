package redis

import (
	"context"
	"encoding/json"
	"time"
)

type CallbackResult[T comparable] struct {
	V   T
	Err error
}

func WaitForKeyCallback[T comparable](ctx context.Context, cl Client, key string) <-chan CallbackResult[T] {
	var v T
	ch := make(chan CallbackResult[T])

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				ch <- CallbackResult[T]{Err: ctx.Err()}
				return
			case <-time.After(time.Millisecond * 100):
				cmd := cl.Get(ctx, key)
				vv, err := cmd.Result()
				if err != nil {
					if err.Error() == "redis: nil" {
						continue
					}
					ch <- CallbackResult[T]{Err: err}
					return
				}
				if vv != "" {
					if err := json.Unmarshal([]byte(vv), &v); err != nil {
						ch <- CallbackResult[T]{Err: err}
						return
					}
					ch <- CallbackResult[T]{V: v}
					return
				}
			}
		}
	}()
	return ch
}

func WaitForKeyCallbackEx[T comparable](ctx context.Context, cl Client, key string, timeout time.Duration) <-chan CallbackResult[T] {
	ctx2, cf := context.WithTimeout(ctx, timeout)
	go func() {
		<-time.After(timeout * 2)
		cf()
	}()
	return WaitForKeyCallback[T](ctx2, cl, key)
}

func NewKeyCallbackWriter[T any](ctx context.Context, cl Client, key string, timeout time.Duration) func(v T) {
	return func(v T) {
		b, err := json.Marshal(v)
		if err != nil {
			return
		}
		cl.Set(ctx, key, b, timeout)
	}
}
