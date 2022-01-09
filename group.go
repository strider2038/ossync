package ossync

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Group struct {
	context context.Context
	cancel  func()

	waiter *sync.WaitGroup

	errOnce *sync.Once
	err     error
}

func NewGroup(ctx context.Context) *Group {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	return &Group{
		cancel:  cancel,
		context: ctx,
		waiter:  &sync.WaitGroup{},
		errOnce: &sync.Once{},
	}
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.waiter.Add(1)

	go func() {
		defer g.waiter.Done()

		err := f(g.context)
		g.errOnce.Do(func() {
			g.err = err
			g.cancel()
		})
	}()
}

func (g *Group) Wait() error {
	g.waiter.Wait()
	g.cancel()

	return g.err
}
