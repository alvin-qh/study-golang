package profile

import (
	"bufio"
	"io"
	"runtime"
	"runtime/pprof"
)

type CpuProfileRecorder struct {
	freq Frequency     // 每次记录的时间间隔
	w    *bufio.Writer // 对文件写入
}

func NewCpuProfileRecorder(w io.Writer, frequency Frequency) *CpuProfileRecorder {
	return &CpuProfileRecorder{
		freq: frequency,
		w:    wrapWriter(w),
	}
}

func (r *CpuProfileRecorder) start() error {
	err := pprof.StartCPUProfile(r.w)
	if err != nil {
		return err
	}

	if r.freq > 0 {
		runtime.SetCPUProfileRate(int(r.freq))
	}

	return nil
}

func (r *CpuProfileRecorder) stop() {
	pprof.StopCPUProfile()
	r.w.Flush()
}
