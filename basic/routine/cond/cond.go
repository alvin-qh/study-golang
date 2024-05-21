package cond

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	condReaderIndex = int32(0) // 计算读协程 ID 的全局变量
)

const (
	COND_TIME_LAYOUT = "2006-01-02T15:04:05.000-07:00"
)

// Cond 对象, 表示一个带 "条件" 的锁 (Condition Lock)
// Cond 对象依赖一个锁对象 (async.Mutex, async.RWMutex), 在其上增加了条件支持
// 在 Cond 对象 Wait 的时候, 会释放锁, 并在 Wait 成功后重新进入锁
// 等待成功意味着有一个 Notify 到达
type FileCond struct {
	filename  string     // 用于读写的文件名
	cond      *sync.Cond // 条件锁对象
	done      bool       // 条件锁等待条件
	notifyAll bool       // 是否通知所有协程
}

// 创建新的 FileCond 对象
func NewFileCond(filename string, notifyAll bool) *FileCond {
	return &FileCond{
		filename:  filename,
		cond:      sync.NewCond(&sync.Mutex{}),
		done:      false,
		notifyAll: notifyAll,
	}
}

// 进入件锁
func (fc *FileCond) lock(name string) {
	fc.cond.L.Lock() // 进入 Cond 对象的锁
	fmt.Printf("%v routine locked\n", name)
}

// 进入条件锁并重设锁条件
func (fc *FileCond) lockAndReset(name string) {
	fc.lock(name)
	fc.done = false // 重设锁条件
}

// 退出条件锁
func (fc *FileCond) unlock(name string) {
	fc.cond.L.Unlock() // 退出 Cond 对象的锁
	fmt.Printf("%v routine unlocked\n", name)
}

// 退出条件锁并, 达成锁条件并发出通知
// 发出通知的方法有两种: `Signal` 和 `Broadcast`, 前者仅通知一个协程, 后者通知所有协程
func (fc *FileCond) unlockAndNotify(name string) {
	defer fc.unlock(name) // 退出锁

	fc.done = true // 设置锁条件达成

	// 发送信号, 表示等待条件已达到,
	if fc.notifyAll {
		fc.cond.Broadcast() // 表示通知到所有的 Wait 函数
	} else {
		fc.cond.Signal() // 表示只通知一个协程, 其它未被通知的协程继续处于 Wait 状态
	}
	fmt.Printf("%v routine broadcast signal\n", name)
}

// 等待锁条件达成
func (fc *FileCond) wait(name string) {
	// 进入等待, 当等待成功后, 再次判断锁条件是否达成
	// 所以一般情况下需要使用条件循环进行处理
	for !fc.done { // 判断读条件是否成立, 如果不成立则进入等待
		fmt.Printf("%v routine begin waiting\n", name)
		fc.cond.Wait()
		fmt.Printf("%v routine wait succeed at %v\n", name, time.Now().Format(COND_TIME_LAYOUT))
	}
}

// 向文件中写入内容, 写入完毕后发送通知
func (fc *FileCond) Write(data []byte) error {
	const name = "[Write]"
	fmt.Printf("%v routine started\n", name)

	fc.lockAndReset(name)          // 锁定并重置锁条件, 此时所有读协程均会进入等待 (Wait 函数)
	defer fc.unlockAndNotify(name) // 处理完毕后解锁并发送通知, 此时等待的协程会结束等待 (Wait 函数)

	// 进入临界区进操作

	f, err := os.Create(fc.filename) // 创建文件用于写操作
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f) // 包装带缓冲区的写对对象
	w.Write(data)           // 写入内容
	w.Flush()

	return nil
}

// 从文件中读取内容, 前提是文件写入结束
func (fc *FileCond) Read() ([]byte, error) {
	name := fmt.Sprintf("[Read(%v)]", atomic.AddInt32(&condReaderIndex, 1))
	fmt.Printf("%v routine started\n", name)

	fc.lock(name)         // 加锁
	defer fc.unlock(name) // 处理完毕后解锁

	fc.wait(name) // 进入等待, 等待 Cond 对象被通知且锁条件达成

	// 加锁成功 (或等待成功), 进入临界区进行读操作

	f, err := os.Open(fc.filename) // 打开文件用于读
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
