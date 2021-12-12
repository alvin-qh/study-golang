package profile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// 内存 Profile 记录结构体
type MemProfileRecorder struct {
	freq Frequency     // 每次记录的时间间隔
	w    *bufio.Writer // 输出记录的 Writer
	ch   chan struct{}
}

// 创建新的内存 Profile
func NewMemProfileRecorder(w io.Writer, frequency Frequency) *MemProfileRecorder {
	return &MemProfileRecorder{
		freq: frequency,
		w:    wrapWriter(w),
		ch:   make(chan struct{}),
	}
}

// 开始记录 MemProfile
func (r *MemProfileRecorder) start() error {
	if r.w != nil {
		go func() {
			for {
				select {
				case <-time.After(time.Millisecond * time.Duration(r.freq)):
					r.record()
				case <-r.ch:
					return
				}
			}
		}()
	}
	return nil
}

// 停止记录 MemProfile
func (r *MemProfileRecorder) stop() {
	if r.ch != nil {
		close(r.ch)
		r.ch = nil
	}

	if r.w != nil {
		r.w.Flush()
		r.w = nil
	}
}

// 记录信息
func (r *MemProfileRecorder) record() {
	n, _ := runtime.MemProfile(nil, true)
	for {
		p := make([]runtime.MemProfileRecord, n+50)
		n, ok := runtime.MemProfile(p, true)
		if ok {
			r.outputMemProfileRecords(p[:n], time.Now().UTC())
			break
		}
	}
}

// 输出 runtime.MemProfileRecord 信息
func (r *MemProfileRecorder) outputMemProfileRecords(records []runtime.MemProfileRecord, recordTime time.Time) {
	if len(records) == 0 {
		return
	}

	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()

	// 输出记录时间
	fmt.Fprintf(r.w, "%v ", recordTime.Format(TIME_LAYOUT_UTC))

	// 计算 records 集合中所有记录的总和
	tp := runtime.MemProfileRecord{}
	for _, record := range records {
		tp.AllocBytes += record.AllocBytes
		tp.AllocObjects += record.AllocObjects
		tp.FreeBytes += record.FreeBytes
		tp.FreeObjects += record.FreeObjects
	}

	// 输出内存使用总和
	fmt.Fprintf(
		r.w, "heap profile: %d: %d [%d: %d] @ heap/%d\n",
		tp.InUseObjects(), tp.InUseBytes(),
		tp.AllocObjects, tp.AllocBytes,
		2*runtime.MemProfileRate,
	)

	// 输出每条内存统计信息及其调用堆栈信息
	for _, record := range records {
		fmt.Fprintf(
			r.w, "%d: %d [%d: %d] @",
			record.InUseObjects(), record.InUseBytes(),
			record.AllocObjects, record.AllocBytes,
		)
		for _, pc := range record.Stack() {
			fmt.Fprintf(r.w, " %#x", pc)
		}
		fmt.Fprintln(r.w)
		r.outputStackRecord(record.Stack(), false)
	}
}

// 输出堆栈信息
func (r *MemProfileRecorder) outputStackRecord(stack []uintptr, allFrames bool) {
	show := allFrames
	wasPanic := false

	// cspell: ignore tracepc
	for i, pc := range stack {
		f := runtime.FuncForPC(pc)
		if f == nil {
			show = true
			fmt.Fprintf(r.w, "#\t%#x\n", pc)
			wasPanic = false
		} else {
			tracepc := pc
			// Back up to call instruction.
			if i > 0 && pc > f.Entry() && !wasPanic {
				if runtime.GOARCH == "386" || runtime.GOARCH == "amd64" {
					tracepc--
				} else {
					tracepc -= 4 // arm, etc
				}
			}
			file, line := f.FileLine(tracepc)
			name := f.Name()
			// Hide runtime.goexit and any runtime functions at the beginning.
			// This is useful mainly for allocation traces.
			wasPanic = name == "runtime.panic"
			if name == "runtime.goexit" || !show && strings.HasPrefix(name, "runtime.") {
				continue
			}
			show = true
			fmt.Fprintf(r.w, "#\t%#x\t%s+%#x\t%s:%d\n", pc, name, pc-f.Entry(), file, line)
		}
	}
	if !show {
		// We didn't print anything; do it again,
		// and this time include runtime functions.
		r.outputStackRecord(stack, true)
		return
	}
	fmt.Fprintln(r.w)
}
