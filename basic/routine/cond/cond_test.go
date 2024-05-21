package cond

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
	COND_FILE_NAME  = "condtest.txt"
	COND_NUM_READER = 3
)

// 测试 FileCond 对象
func TestCondBroadcasting(t *testing.T) {
	defer os.Remove(COND_FILE_NAME)

	wg := sync.WaitGroup{}
	wg.Add(COND_NUM_READER)

	fc := NewFileCond(COND_FILE_NAME, true)

	// 进行 3 次读操作, 因为写操作尚未进行, 所以读操作无法进行
	for i := 0; i < COND_NUM_READER; i++ {
		go func() {
			defer wg.Done()

			data, err := fc.Read()
			assert.Nil(t, err)
			assert.Equal(t, "Hello World!", string(data))
		}()
	}

	// 稍等后进入写操作, 以保证读操作均进入 Wait 操作
	time.Sleep(100 * time.Millisecond)

	go func() {
		err := fc.Write([]byte("Hello World!"))
		assert.Nil(t, err)
	}()

	wg.Wait()
}

func TestCondSignaling(t *testing.T) {
	defer os.Remove(COND_FILE_NAME)

	wg := sync.WaitGroup{}
	wg.Add(COND_NUM_READER)

	fc := NewFileCond(COND_FILE_NAME, false)

	// 进行 3 次读操作, 因为写操作尚未进行, 所以读操作无法进行
	// 每次打开文件从头读取, 读取最近一次写入文件的内容
	for i := 0; i < COND_NUM_READER; i++ {
		go func() {
			defer wg.Done()

			data, err := fc.Read()
			assert.Nil(t, err)

			fmt.Printf("[Result] Read content is: %v\n", string(data))
		}()
	}

	// 稍等后进入写操作, 以保证读操作均进入 Wait 操作
	time.Sleep(100 * time.Millisecond)

	// 进行 3 次写操作, 每次写操作通知一个协程处理对应的读操作
	// 注意, 这三次写操作并不是追加写, 而是每次清空文件进行写, 所以操作结束后, 文件中只有最后一次写的内容
	for i := 0; i < COND_NUM_READER; i++ {
		num := i + 1
		go func() {
			err := fc.Write([]byte(fmt.Sprintf("[%v] Hello World!\n", num)))
			assert.Nil(t, err)
		}()
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
}
