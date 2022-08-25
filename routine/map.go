package routine

import (
	"context"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gitlab.momoso.com/cm/kit/third_party/lg"
	"golang.org/x/sync/errgroup"
)

// Map is a helper function to parallel operate on slice.
// Example usage:
// ```
//   itemIDs := []int{1, 2, 3}
//   items := make([]*Item, len(itemIDs))
//   routine.Map(len(items), 10, func(i int) {
//     items[i] = getItemByID(itemIDs[i])
//   })
// ```
func Map(size, numWorker int, f func(i int) error) error {
	return MapWithTimeout(size, numWorker, time.Hour*24*365, func(ctx context.Context, i int) error {
		err := f(i)
		return err
	})
}

// MapWithTimeout is the same as Map, except each function need to be executed within timeout.
// Important reminder to handle the context provided:
// It is better to check the `ctx.Err() == nil` before any write operation, because the
// function may be canceled at any momement, and the underlying resource refering may have already been
// discarded.
func MapWithTimeout(size, numWorker int, timeout time.Duration, f func(ctx context.Context, i int) error) error {
	return MapWithCtxTimeout(context.Background(), size, numWorker, timeout, f)
}

func MapWithContext(ctx context.Context, size, numWorker int, f func(ctx context.Context, i int) error) error {
	return MapWithCtxTimeout(ctx, size, numWorker, time.Hour*24*365, f)
}

func MapWithCtxTimeout(baseCtx context.Context, size, numWorker int, timeout time.Duration, f func(ctx context.Context, i int) error) error {
	if numWorker <= 0 {
		return errors.New("no num worker has been specified")
	}
	if numWorker == 1 {
		return doSimpleMap(baseCtx, size, f)
	}
	var allerror error
	ch := make(chan int)
	var mutex sync.Mutex
	go func() {
		defer close(ch)
		for i := 0; i < size; i++ {
			ch <- i
		}
	}()

	delayTerminate := false
	grp := errgroup.Group{}
	for i := 0; i < numWorker; i++ {
		grp.Go((func() error {
			for idx := range ch {
				if baseCtx.Err() != nil {
					mutex.Lock()
					allerror = multierror.Append(allerror, errors.Wrapf(baseCtx.Err(), "Task[%d]", idx))
					mutex.Unlock()
					break
				}
				ctx, cancel := context.WithTimeout(baseCtx, timeout)
				start := time.Now()
				err := f(ctx, idx)
				if err != nil {
					mutex.Lock()
					allerror = multierror.Append(allerror, errors.Wrapf(baseCtx.Err(), "Task[%d]", idx))
					mutex.Unlock()
				}
				if ctx.Err() != nil && time.Since(start) > timeout+time.Second {
					delayTerminate = true
				}
				if ctx.Err() != nil {
					mutex.Lock()
					allerror = multierror.Append(allerror, errors.Wrapf(baseCtx.Err(), "Task[%d]", idx))
					mutex.Unlock()

				}
				cancel()
			}
			return nil
		}))
	}
	grp.Wait()

	for range ch {
		//pass
	}
	if delayTerminate {
		lg.Warn("worker dont respect the ctx", lg.GetFuncName(f))
	}
	return allerror
}

func doSimpleMap(ctx context.Context, size int, f func(ctx context.Context, i int) error) error {
	var allerror error
	for i := 0; i < size; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := f(ctx, i)
		if err != nil {
			allerror = multierror.Append(allerror, err)
		}
	}
	return allerror
}
