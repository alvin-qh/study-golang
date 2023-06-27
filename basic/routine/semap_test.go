package routine

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

const (
	SEMAP_FILE_NAME  = "semaphore.txt"
	SEMAP_NUM_READER = 3
)

// 测试 FileCond 对象
func TestSemaphoreReadWrite(t *testing.T) {
	defer os.Remove(SEMAP_FILE_NAME)

	wg := sync.WaitGroup{}
	wg.Add(COND_NUM_READER)

	fs, err := NewFileSemaphore(SEMAP_FILE_NAME, SEMAP_NUM_READER)
	assert.NoError(t, err)

	// 进行 3 次读操作, 因为写操作尚未进行, 所以读操作无法进行
	for i := 0; i < SEMAP_NUM_READER; i++ {
		go func() {
			defer wg.Done()

			data, err := fs.Read()
			assert.NoError(t, err)

			fmt.Printf("[Result] Read content is: %v\n", string(data))
		}()
	}

	go func() {
		err := fs.Write([]byte("[1-2] Hello World!"), 2)
		assert.NoError(t, err)
	}()

	// 第一次写入完毕后稍等片刻, 让对应的读操作结束, 否则两次连续的写操作会导致内容混乱
	time.Sleep(time.Second)

	go func() {
		err := fs.Write([]byte("[3] Hello World!"), 1)
		assert.NoError(t, err)
	}()

	wg.Wait()
}
