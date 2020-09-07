package flexibale

import (
	"context"
	"errors"
	"github/WowVeryLogin/multer/pkg/pool"
	"sync/atomic"
	"testing"
)

func TestWork(t *testing.T) {
	t.Parallel()
	p := New(&Config{
		MaxWorkers: 2,
	})
	defer p.Close()

	var totalSum int64

	p.Put(context.Background(), pool.Task(func(ctx context.Context) error {
		atomic.AddInt64(&totalSum, 1)
		return nil
	}))
	p.Put(context.Background(), pool.Task(func(ctx context.Context) error {
		atomic.AddInt64(&totalSum, 1)
		return nil
	}))
	p.Put(context.Background(), pool.Task(func(ctx context.Context) error {
		atomic.AddInt64(&totalSum, 1)
		return nil
	}))
	p.WaitBarrier(nil)
	if totalSum != 3 {
		t.FailNow()
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	p := New(&Config{
		MaxWorkers: 2,
	})
	defer p.Close()

	var totalSum int64

	p.Put(context.Background(), pool.Task(func(ctx context.Context) error {
		return errors.New("some error")
	}))
	p.Put(context.Background(), pool.Task(func(ctx context.Context) error {
		atomic.AddInt64(&totalSum, 1)
		return nil
	}))
	p.Put(context.Background(), pool.Task(func(ctx context.Context) error {
		atomic.AddInt64(&totalSum, 1)
		return nil
	}))

	var err error
	p.WaitBarrier(func(e error) {
		err = e
	})
	if err == nil {
		t.FailNow()
	}
}
