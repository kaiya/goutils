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

func MapWithCtxTimeout(baseCtx context.Context, size, numWorker int, timeout time.Duration, f func(ctx context.Context, i int) error) error {
	if numWorker <= 0 {
		return errors.New("no num worker has been specified")
	}
	if numWorker == 1 {
		// do simple map
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
