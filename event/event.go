package event

import (
	"context"
	"strings"
	"sync"

	"github.com/fengjx/go-halo/halo"
)

// Event 事件定义
type Event[T any] string

// HandlerFunc 事件处理函数
type HandlerFunc[T any] func(data T)

var (
	handlerMap = make(map[string][]any)

	handlerLock sync.Mutex
)

// On Register a handler for an event.
func On[T any](e Event[T], h HandlerFunc[T]) {
	name := strings.TrimSpace(string(e))
	if name == "" {
		panic("event name cannot be empty")
	}
	handlerLock.Lock()
	defer handlerLock.Unlock()
	if hs, ok := handlerMap[name]; ok {
		handlerMap[name] = append(hs, h)
	} else {
		handlerMap[name] = []any{h}
	}
}

// Emit Trigger an event.
func Emit[T any](e Event[T], data T) {
	for _, h := range handlerMap[string(e)] {
		fn := h.(HandlerFunc[T])
		halo.GracefulRun(func(ctx context.Context) {
			fn(data)
		})
	}
}
