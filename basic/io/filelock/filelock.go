//go:build !windows

package filelock

import (
	"os"
	"runtime"
	"syscall"
)

// 和文件锁相关的常量包括:
//
//	LOCK_EX: 互斥锁
//	LOCK_NB: 非阻塞
//	LOCK_SH: 共享锁
//	LOCK_UN: 解锁
//
// 文件锁结构体
type FileLock struct {
	dir      string   // 加锁使用的文件路径
	f        *os.File // 锁文件对象
	nonBlock bool     // 是否阻塞
}

// 新建一个文件锁对象
func New(dir string, nonBlock bool) *FileLock {
	// 设置锁定文件路径
	fl := &FileLock{dir: dir, nonBlock: nonBlock}

	// 在引用失效后, 自动解锁
	runtime.SetFinalizer(fl, func(fl *FileLock) { fl.Unlock() })
	return fl
}

// 关闭文件锁
func (fl *FileLock) Close() error {
	if fl.f != nil {
		if err := fl.f.Close(); err != nil {
			return err
		}
	}

	if err := os.Remove(fl.dir); err != nil {
		return err
	}
	return nil
}

// 加锁操作, 在文件上施加一个互斥锁
func (fl *FileLock) XLock() error {
	// 创建或打开加锁文件
	f, err := os.OpenFile(fl.dir, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// 设置锁标志
	// 设置互斥锁
	how := syscall.LOCK_EX
	if fl.nonBlock {
		// 设置是否阻塞
		how |= syscall.LOCK_NB
	}

	// 对文件进行加锁, Fd() 函数返回文件描述符句柄
	if err = syscall.Flock(int(f.Fd()), how); err != nil {
		f.Close()
		return err
	}

	// 保持文件对象
	fl.f = f
	return nil
}

// 判断文件是否已被锁定
func (fl *FileLock) IsLocked() bool { return fl.f != nil }

// 解锁
func (fl *FileLock) Unlock() error {
	if fl.f == nil {
		return nil
	}

	// 关闭文件对象
	defer func() {
		fl.f.Close()
		fl.f = nil
	}()

	// 解锁
	return syscall.Flock(int(fl.f.Fd()), syscall.LOCK_UN)
}
