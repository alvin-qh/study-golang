package routine

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

func TestChanLock(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	ctx := context.TODO()

	l := NewLock()

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()

			locked := l.Lock(ctx)
			assert.True(t, locked)
			if locked {
				defer l.Unlock()
				time.Sleep(time.Second)
			}
		}()
	}

	wg.Wait()
}
