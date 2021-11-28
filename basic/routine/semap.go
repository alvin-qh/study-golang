package routine

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

// 创建新的 FileSemaphore 对象
// semaphore.Weighted 结构体表示一个信号量
// 创建 Weighted 结构体时，需设置信号量的总数，使用 Release 函数释放一个信号，使用 Acquire 函数等待被释放的信号量并重新占用它
//
// 信号量需要使用到 Google 的扩展库，即 go get golang.org/x/sync/semaphore
type FileSemaphore struct {
	filename string              // 用于读写的文件名
	weighted *semaphore.Weighted // 信号量
	ctx      context.Context     // 信号量上下文
}

var (
	semapWriterIndex = int32(0) // 计算写协程 ID 的全局变量
	semapReaderIndex = int32(0) // 计算读协程 ID 的全局变量
)

// 创建新的 FileSemaphore 对象，count 表示信号量的整体数量
func NewFileSemaphore(filename string, count int64) (*FileSemaphore, error) {
	fs := &FileSemaphore{
		filename: filename,                     // 用于操作的文件名
		weighted: semaphore.NewWeighted(count), // 信号量，设置其总共可使用的数量，此时所有信号量均是占用状态
		ctx:      context.Background(),         // 信号量上下文，参考 routine/context.go
	}
	fs.acquire(count, "[Init]") // 将所有的信号量设置为“已占用”状态，等待释放
	return fs, nil
}

// 占用信号
// 只有释放一个信号量，才能成功占用一个信号量，否则会进行等待
// 也可也使用 TryAcquire 进行尝试占用，这个方法不会阻塞，而是返回 bool 值表示是否有可用的信号量
func (fs *FileSemaphore) acquire(n int64, name string) error {
	fmt.Printf("%v routine, %v semaphores are acquiring\n", name, n)

	err := fs.weighted.Acquire(fs.ctx, n)
	if err != nil {
		return err
	}
	fmt.Printf("%v routine, %v semaphores was acquired\n", name, n)
	return nil
}

// 释放信号
// 当有被占用的信号量时，可以释放指定数量的被占用信号量，此时等待信号量的协程可以等待成功，并继续执行
func (fs *FileSemaphore) release(n int64, name string) {
	fs.weighted.Release(n)
	fmt.Printf("%v routine, %v semaphores was released\n", name, n)
}

// 向文件中写入内容，写入完毕后发送信号
func (fs *FileSemaphore) Write(data []byte, releaseCount int64) error {
	name := fmt.Sprintf("[Write(%v)]", atomic.AddInt32(&semapWriterIndex, 1))
	fmt.Printf("%v routine started\n", name)

	defer fs.release(releaseCount, name) // 处理完毕后发送信号，释放 releaseCount 个协程

	// 进入临界区进操作

	f, err := os.Create(fs.filename) // 创建文件用于写操作
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f) // 包装带缓冲区的写对对象
	w.Write(data)           // 写入内容
	w.Flush()

	return nil
}

// 等待信号，并读取文件
func (fs *FileSemaphore) Read() ([]byte, error) {
	name := fmt.Sprintf("[Read(%v)]", atomic.AddInt32(&semapReaderIndex, 1))
	fmt.Printf("%v routine started\n", name)

	if err := fs.acquire(1, name); err != nil { // 等待有信号被释放，并占用该信号
		return nil, err
	}

	// 信号占用成功，进入临界区进行读操作

	f, err := os.Open(fs.filename) // 打开文件用于读
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)      // 为文件包装带缓冲区的读对象
	data, _, err := r.ReadLine() // 读取一行文本
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v routine read succeed\n", name)
	return data, nil
}
