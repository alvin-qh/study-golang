package routine

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

var (
	semapWriterIndex = int32(0) // 计算写协程 ID 的全局变量
	semapReaderIndex = int32(0) // 计算读协程 ID 的全局变量
)

// 创建新的 FileSemaphore 对象
// semaphore.Weighted 结构体表示一个信号量
// 创建 Weighted 结构体时, 需设置信号量的总数, 使用 Release 函数释放一个信号, 使用 Acquire 函数等待被释放的信号量并重新占用它
//
// 信号量需要使用到 Google 的扩展库, 即 go get golang.org/x/sync/semaphore
type FileSemaphore struct {
	filename string              // 用于读写的文件名
	weighted *semaphore.Weighted // 信号量
	ctx      context.Context     // 信号量上下文
}

// 创建 `FileSemaphore` 实例
//
// `count` 表示信号量的整体数量
func New(filename string, count int64) (*FileSemaphore, error) {
	fs := &FileSemaphore{
		filename: filename,                     // 用于操作的文件名
		weighted: semaphore.NewWeighted(count), // 信号量, 设置其总共可使用的数量, 此时所有信号量均是占用状态
		ctx:      context.Background(),         // 信号量上下文, 参考 routine/context.go
	}
	fs.acquire(count, "[Init]") // 将所有的信号量设置为 "已占用" 状态, 等待释放
	return fs, nil
}

// 占用信号量
//
// 只有一方释放一个信号量, 另一方才能成功占用一个信号量, 否则后者会进行等待
//
// 也可也使用 `semaphore.Weighted` 实例的 `TryAcquire` 方法进行尝试占用,
// 这个方法不会阻塞, 而是返回 `bool` 值表示是否有可用的信号量
func (fs *FileSemaphore) acquire(n int64, name string) error {
	fmt.Printf("%v routine, %v semaphores are acquiring\n", name, n)

	if err := fs.weighted.Acquire(fs.ctx, n); err != nil {
		return err
	}
	fmt.Printf("%v routine, %v semaphores was acquired\n", name, n)
	return nil
}

// 释放信号量
//
// 当有被占用的信号量时, 可以释放指定数量的被占用信号量, 此时等待信号量的一方可以等待成功并继续执行
func (fs *FileSemaphore) release(n int64, name string) {
	fs.weighted.Release(n)
	fmt.Printf("%v routine, %v semaphores was released\n", name, n)
}

// 向文件中写入内容
//
// 写入完毕后发送信号, 等待写入完毕的一方可以等待成功
func (fs *FileSemaphore) Write(data []byte, releaseCount int64) error {
	name := fmt.Sprintf("[Write(%v)]", atomic.AddInt32(&semapWriterIndex, 1))
	fmt.Printf("%v routine started\n", name)

	// 处理完毕后发送信号, 释放 releaseCount 个协程
	defer fs.release(releaseCount, name)

	// 进入临界区进操作

	// 创建文件用于写操作
	f, err := os.Create(fs.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// 包装带缓冲区的写对对象
	w := bufio.NewWriter(f)

	// 写入内容
	w.Write(data)
	w.Flush()

	return nil
}

// 等待信号, 并读取文件
func (fs *FileSemaphore) Read() ([]byte, error) {
	name := fmt.Sprintf("[Read(%v)]", atomic.AddInt32(&semapReaderIndex, 1))
	fmt.Printf("%v routine started\n", name)

	// 等待有信号被释放, 并占用该信号
	if err := fs.acquire(1, name); err != nil {
		return nil, err
	}

	// 信号占用成功, 进入临界区进行读操作

	// 打开文件用于读
	f, err := os.Open(fs.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 为文件包装带缓冲区的读对象
	r := bufio.NewReader(f)

	// 读取一行文本
	data, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v routine read succeed\n", name)
	return data, nil
}
