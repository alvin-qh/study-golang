package memory

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	// 以 Z 结尾的标准 UTC 时间格式
	TIME_LAYOUT_UTC = "2006-01-02T15:04:05.000Z"
)

type Frequency time.Duration

// 记录内存 Profile 信息的结构体
type Profile struct {
	freq Frequency     // 每次记录的时间间隔
	w    *bufio.Writer // 输出记录的 Writer
	ch   chan struct{}
}

// 将 `io.Writer` 包装为 `bufio.Writer`
func wrapWriter(w io.Writer) *bufio.Writer {
	if bufW, ok := w.(*bufio.Writer); ok {
		return bufW
	}
	return bufio.NewWriter(w)
}

// 创建 `Profile` 类型实例
func New(w io.Writer, frequency Frequency) *Profile {
	return &Profile{
		freq: frequency,
		w:    wrapWriter(w),
		ch:   make(chan struct{}),
	}
}

// 向 `io.Writer` 输出格式化后的字符串
func (r *Profile) printf(format string, a ...any) {
	if r.w != nil {
		if _, err := fmt.Fprintf(r.w, format, a...); err != nil {
			panic(err)
		}
	}
}

// 向 `io.Writer` 输出字符串和换行
func (r *Profile) println(a ...any) {
	if r.w != nil {
		if _, err := fmt.Fprintln(r.w, a...); err != nil {
			panic(err)
		}
	}
}

// 输出堆栈帧信息
//
// 堆栈帧为当前调用位置的程序计数器 (`pc`), 通过该计数器可以找到该帧所表示的调用位置的信息, 包括:
//   - 文件名
//   - 行号
//   - 函数名
func (r *Profile) outputStackFrame(stack []uintptr, allFrames bool) {
	show := allFrames
	wasPanic := false

	for i, pc := range stack {
		// 通过程序计数器获取函数调用位置
		f := runtime.FuncForPC(pc)
		if f == nil {
			show = true
			r.printf("#\t%#x\n", pc)
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
			r.printf("#\t%#x\t%s+%#x\t%s:%d\n", pc, name, pc-f.Entry(), file, line)
		}
	}

	if !show {
		// We didn't print anything; do it again,
		// and this time include runtime functions.
		r.outputStackFrame(stack, true)
		return
	}
	r.println()
}

// 输出内存 Profile 记录数据
func (r *Profile) outputMemRecords(records []runtime.MemProfileRecord, recordTime time.Time) {
	if len(records) == 0 {
		return
	}

	// 对出现的错误进行恢复
	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()

	// 输出记录时间
	r.printf("%v ", recordTime.Format(TIME_LAYOUT_UTC))

	// 计算 records 集合中所有记录的总和
	tp := runtime.MemProfileRecord{}
	for _, record := range records {
		tp.AllocBytes += record.AllocBytes
		tp.AllocObjects += record.AllocObjects
		tp.FreeBytes += record.FreeBytes
		tp.FreeObjects += record.FreeObjects
	}

	// 输出内存使用总和
	r.printf(
		"heap profile: %d: %d [%d: %d] @ heap/%d\n",
		tp.InUseObjects(), tp.InUseBytes(),
		tp.AllocObjects, tp.AllocBytes,
		2*runtime.MemProfileRate,
	)

	// 输出每条内存统计信息及其调用堆栈信息
	for _, record := range records {
		r.printf(
			"%d: %d [%d: %d] @",
			record.InUseObjects(), record.InUseBytes(),
			record.AllocObjects, record.AllocBytes,
		)
		for _, pc := range record.Stack() {
			r.printf(" %#x", pc)
		}
		r.println()
		r.outputStackFrame(record.Stack(), false)
	}
}

// 读取内存中记录的 Profile 数据
//
// 通过 `runtime.MemProfile` 函数可以读取内存中记录的 Profile 数据 (读取并删除), 其中:
//   - 第一个参数为存储 Profile 的集合实例, 如果该集合长度小于实际的 Profile 数量, 则不写入数据, 但返回实际 Profile 数量
//   - 第二个参数为是否返回已经读取的 Profile 数据
func (r *Profile) record() {
	// 获取已经积累的 Profile 数量
	n, _ := runtime.MemProfile(nil, true)
	for {
		// 分配内存读取 Profile 数据, 为防止读取失败, 需要冗余分配一部分空间
		p := make([]runtime.MemProfileRecord, n+10)

		// 读取 Profile 数据
		n, ok := runtime.MemProfile(p, true)
		if ok {
			// 将读取的 Profile 数据输出到 `io.Writer` 实例中
			r.outputMemRecords(p[:n], time.Now().UTC())
			break
		}
	}
}

// 开始记录内存 Profile
//
// 通过一个定时器, 在指定的间隔时间后向 `io.Writer` 中写入内存 Profile 数据
func (r *Profile) Start() error {
	if r.w != nil {
		// 启动异步函数, 在定时器时间到达时进行一次内存 Profile 数据的记录
		go func() {
			for {
				select {
				case <-time.After(time.Millisecond * time.Duration(r.freq)): // 等待定时器时间到达
					r.record() // 记录内存 Profile 数据
				case <-r.ch: // 通道被关闭, 退出循环, 结束协程
					if r.w != nil {
						r.w.Flush()
						r.w = nil
					}
					r.ch = nil
					return
				}
			}
		}()
	}
	return nil
}

// 停止记录内存 Profile
func (r *Profile) Stop() {
	if r.ch != nil {
		close(r.ch)
	}
}
