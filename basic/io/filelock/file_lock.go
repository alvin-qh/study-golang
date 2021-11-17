package filelock

import (
	"os"
	"runtime"
	"syscall"
)

type FileLock struct {
	dir string
	f   *os.File
}

func New(dir string) *FileLock {
	fl := &FileLock{dir: dir}
	runtime.SetFinalizer(fl, func(fl *FileLock) { fl.Close() })
	return fl
}

func (fl *FileLock) Close() error {
	if fl.f == nil {
		return nil
	}
	return fl.f.Close()
}

func (fl *FileLock) Lock(block bool) error {
	f, err := os.OpenFile(fl.dir, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	how := syscall.LOCK_EX
	if block {
		how |= syscall.LOCK_NB
	}
	return syscall.Flock(int(f.Fd()), how)
}

func (fl *FileLock) Unlock() error {
	if fl.f == nil {
		return nil
	}
	defer func() {
		fl.f.Close()
		fl.f = nil
	}()

	return syscall.Flock(int(fl.f.Fd()), syscall.LOCK_UN)
}
