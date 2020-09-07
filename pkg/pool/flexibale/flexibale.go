package flexibale

import (
	"context"
	"github/WowVeryLogin/multer/pkg/pool"
	"sync"
	"sync/atomic"
)

type Config struct {
	MaxWorkers int `json:"max_workers"`
}

type poolImpl struct {
	*Config
	tasks      chan func() error
	numWorkers int64
	errors     chan error

	barrier chan struct{}
	wg      sync.WaitGroup
}

func New(cfg *Config) pool.Pool {
	return &poolImpl{
		Config:  cfg,
		tasks:   make(chan func() error),
		errors:  make(chan error, cfg.MaxWorkers),
		barrier: make(chan struct{}),
	}
}

func (p *poolImpl) Errors() <-chan error {
	return p.errors
}

func (p *poolImpl) WaitBarrier(onerror func(err error)) {
	go func() {
		p.wg.Wait()
		p.barrier <- struct{}{}
	}()

	for {
		select {
		case <-p.barrier:
			return
		case err := <-p.errors:
			if onerror != nil {
				onerror(err)
			}
		}
	}
}

func (p *poolImpl) Close() {
	close(p.tasks)
	close(p.barrier)
	close(p.errors)
}

func (p *poolImpl) workerFn() {
	for t := range p.tasks {
		err := t()
		if err != nil {
			p.errors <- err
		}
		p.wg.Done()
	}
}

func (p *poolImpl) Put(ctx context.Context, task pool.Task) {
	p.wg.Add(1)
	t := func() error {
		select {
		case <-ctx.Done():
			return nil
		default:
			return task(ctx)
		}
	}
	select {
	case p.tasks <- t:
		return
	default:
		numWorkers := atomic.AddInt64(&p.numWorkers, 1)
		if numWorkers <= int64(p.MaxWorkers) {
			go p.workerFn()
		}
	}

	p.tasks <- t
}
