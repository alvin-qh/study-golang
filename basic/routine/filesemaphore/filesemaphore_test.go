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

// 测试 `FileSemaphore`
//
// 通过信号量, 在一方完成文件写入完毕后, 等待的其它方得到信号, 开始读取文件
func TestSemaphoreReadWrite(t *testing.T) {
	// 实例化等待组, 共对 3 个任务进行等待
	wg := sync.WaitGroup{}
	wg.Add(3)

	// 实例化带信号量的文件, 共指定 3 个信号量
	fs, err := New(SEMAP_FILE_NAME, 3)
	assert.Nil(t, err)

	defer os.Remove(SEMAP_FILE_NAME)

	// 启动三个 routine 对文件进行读操作
	// 每完成一次写操作, 释放一个信号量, 这里的一个读操作信号量可以被释放, 完成一个 routine
	for i := 0; i < SEMAP_NUM_READER; i++ {
		go func() {
			defer wg.Done()

			data, err := fs.Read()
			assert.Nil(t, err)

			fmt.Printf("[Result] Read content is: %v\n", string(data))
		}()
	}

	// 启动一个 routine 对文件进行写操作, 并释放 2 个信号量, 所以当本次写操作完成后, 可以释放 2 次读操作
	go func() {
		err := fs.Write([]byte("[1-2] Hello World!"), 2)
		assert.Nil(t, err)
	}()

	// 写入完毕后稍等片刻, 让对应的读操作结束, 否则两次连续的写操作会导致内容混乱
	time.Sleep(time.Millisecond * 200)

	// 再次进行一次写操作, 并释放 1 个信号量
	go func() {
		err := fs.Write([]byte("[3] Hello World!"), 1)
		assert.Nil(t, err)
	}()

	// 等待所有 routine 结束
	wg.Wait()
}
