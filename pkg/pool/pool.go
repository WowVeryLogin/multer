package pool

import "context"

type Task func(context.Context) error

type Pool interface {
	Close()
	Put(context.Context, Task)
	WaitBarrier(func(error))
	Errors() <-chan error
}
