package main

import (
	"context"
	"sync"
)

type Limiter struct {
	limit chan struct{}
	wg    sync.WaitGroup
}

func NewLimiter(n int) *Limiter {
	return &Limiter{limit: make(chan struct{}, n)}
}

func (lim *Limiter) Go(ctx context.Context, id int, fn func(id int)) bool {
	if ctx.Err() != nil {
		return false
	}

	select {
	case lim.limit <- struct{}{}:
	case <-ctx.Done():
		return false
	}

	lim.wg.Add(1)
	go func() {
		defer func() {
			<-lim.limit
			lim.wg.Done()
		}()

		fn(id)
	}()

	return true
}

func (lim *Limiter) Wait() {
	lim.wg.Wait()
}
