package record

import (
	"io"
	"runtime/pprof"
)

// 定义记录各类 Profile 信息的结构体
//
// 记录自记录开始后程序对 CPU 的使用情况
// 记录的信息以二进制形式保存在指定的 `Writer` 对象 (例如一个文件) 中, 通过 Go 语言提供的工具可以对其进行分析
//
//	go tool pprof ./path/cpu.profile
//
// 上述命令打开指定的 Profile 文件, 之后即可通过命令显示各类分析结果, 例如:
//   - `top <N>`      显示前 N 个 CPU 使用率最高的函数调用
//   - `web`          生成一个 svg 图形文件, 并通过默认浏览器打开查看
//   - `list <regex>` 根据所给的正则表达式, 列出匹配的采样位置附近的代码
//   - `help`         显示帮助信息
//
// 以图像界面打开分析文件
//
//	go tool pprof -http :8080 ./path/cpu.profile
//
// 上述命令在 `8080` 端口启动一个 Web Server, 可以以图形化的方式查看各种 CPU 使用数据
//
// 火焰图: 要以火焰图展示 CPU 使用情况, 需要安装额外的组件
//
//	sudo apt install -y graphviz
//
// `top` 命令列出如下内容:
//
//	Showing nodes accounting for 1.93s, 96.02% of 2.01s total
//	Dropped 49 nodes (cum <= 0.01s)
//	Showing top 10 nodes out of 35
//	  flat  flat%   sum%        cum   cum%
//	 0.89s 44.28% 44.28%      0.89s 44.28%  runtime.memmove
//	 0.57s 28.36% 72.64%      0.57s 28.36%  runtime.nanotime (inline)
//	 0.24s 11.94% 84.58%      0.24s 11.94%  runtime.memclrNoHeapPointers
//	 .....
//
// 其中: `flat` 表示某个函数执行占用的 CPU 时间和百分比; `sum` 表示当前 `flat%` 和之前 `flat%` 的累加值; `cum` 表示某个函数及其子函数运行占用的 CPU 时间和百分比
//
// 在测试中使用 Profile 记录性能数据:
//   - 记录 CPU 使用情况: `go test ./runtime/profile -test.cpuprofile cpu.profile`  # 记录测试的 CPU 使用情况, 保存在 `cpu.profile` 文件中
//   - 记录 MEM 使用情况: `go test ./runtime/profile -test.memprofile mem.profile`  # 记录测试的 内存 使用情况, 保存在 `mem.profile` 文件中
//   - 记录文件可以通过: `go tool pprof cpu.profile` 或者 `go tool pprof mem.profile` 查看
type Recorder struct {
	cpuW  io.Writer
	heapW io.Writer
}

// 创建实例
func New(cpuW io.Writer, heapW io.Writer) *Recorder {
	return &Recorder{
		cpuW, heapW,
	}
}

// 开始记录 CPU Profile
func (r *Recorder) Start() error {
	if r.cpuW != nil {
		if err := pprof.StartCPUProfile(r.cpuW); err != nil {
			return err
		}
	}
	return nil
}

// 结束记录 CPU Profile, 并记录 Heap Profile
func (r *Recorder) Stop() error {
	if r.cpuW != nil {
		pprof.StopCPUProfile()
	}

	if r.heapW != nil {
		if err := pprof.WriteHeapProfile(r.heapW); err != nil {
			return err
		}
	}
	return nil
}
