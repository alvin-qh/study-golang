package filelock

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const LOCK_FILE = "./.lock"

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestFileLockAndUnlock(t *testing.T) {
	// defer os.Remove(LOCK_FILE)

	tasks := map[string]bool{
		"task1": false,
		"task2": false,
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	worker := func(name string) {
		if _, exist := tasks[name]; exist {
			tasks[name] = true
			wg.Done()
		}
	}

	go worker("task1")
	go worker("task2")

	waitTimeout(wg, time.Second*2)
	// assert.False(t, ok)

	fl := New(LOCK_FILE)
	err := fl.Unlock()
	assert.NoError(t, err)

	waitTimeout(wg, time.Second*2)
	// assert.True(t, ok)

	assert.True(t, tasks["task1"])
	assert.True(t, tasks["task2"])
}
