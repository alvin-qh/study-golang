package profile

import (
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/**
 * 记录 Profile
 *
 * 记录自记录开始后程序对 CPU 的使用情况
 * 记录的信息以二进制形式保存在指定的 `Writer` 对象 (例如一个文件) 中, 通过 Go 语言提供的工具可以对其进行分析
 *
 *   ```bash
 *   go tool pprof ./path/cpu.profile
 *   ```
 *
 * 上述命令打开指定的 Profile 文件, 之后即可通过命令显示各类分析结果, 例如:
 *   `top <N>`      显示前 N 个 CPU 使用率最高的函数调用
 *   `web`          生成一个 svg 图形文件, 并通过默认浏览器打开查看
 *   `list <regex>` 根据所给的正则表达式, 列出匹配的采样位置附近的代码
 *   `help`         显示帮助信息
 *
 * 以图像界面打开分析文件
 *
 *   ```bash
 *   go tool pprof -http :8080 ./path/cpu.profile
 *   ```
 *
 * 上述命令在 `8080` 端口启动一个 Web Server, 可以以图形化的方式查看各种 CPU 使用数据
 *
 * 火焰图: 要以火焰图展示 CPU 使用情况, 需要安装额外的组件
 *
 *   ```bash
 *   sudo apt install -y graphviz
 *   ```
 *
 * `top` 命令列出如下内容:
 *    Showing nodes accounting for 1.93s, 96.02% of 2.01s total
 *    Dropped 49 nodes (cum <= 0.01s)
 *    Showing top 10 nodes out of 35
 *         flat  flat%   sum%        cum   cum%
 *        0.89s 44.28% 44.28%      0.89s 44.28%  runtime.memmove
 *        0.57s 28.36% 72.64%      0.57s 28.36%  runtime.nanotime (inline)
 *        0.24s 11.94% 84.58%      0.24s 11.94%  runtime.memclrNoHeapPointers
 *        .....
 *
 * 其中: `flat` 表示某个函数执行占用的 CPU 时间和百分比; `sum` 表示当前 `flat%` 和之前 `flat%` 的累加值; `cum` 表示某个函数及其子函数运行占用的 CPU 时间和百分比
 *
 * 在测试中使用 Profile 记录性能数据:
 *   记录 CPU 使用情况: `go test ./runtime/profile -test.cpuprofile cpu.profile`  # 记录测试的 CPU 使用情况, 保存在 `cpu.profile` 文件中
 *   记录 MEM 使用情况: `go test ./runtime/profile -test.memprofile mem.profile`  # 记录测试的 内存 使用情况, 保存在 `mem.profile` 文件中
 *   记录文件可以通过: `go tool pprof cpu.profile` 或者 `go tool pprof mem.profile` 查看
 */

const (
	MEM_PROFILE_FILENAME  = "mem.profile"
	CPU_PROFILE_FILENAME  = "cpu.profile"
	HEAP_PROFILE_FILENAME = "heap.profile"
)

// 测试记录 Profile 数据
func TestRecordProfile(t *testing.T) {
	defer os.Remove(MEM_PROFILE_FILENAME)
	defer os.Remove(CPU_PROFILE_FILENAME)
	defer os.Remove(HEAP_PROFILE_FILENAME)

	memf, err := os.Create(MEM_PROFILE_FILENAME) // 创建记录 NewMemProfile 信息的文件
	assert.NoError(t, err)
	defer memf.Close()

	profile := NewMemProfile(memf, 500) // 创建 NewMemProfile 对象

	cpuf, err := os.Create(CPU_PROFILE_FILENAME)
	assert.NoError(t, err)
	defer cpuf.Close()

	err = profile.Start() // 开始记录内存使用信息
	assert.NoError(t, err)
	defer profile.Stop() // 函数结束后停止记录内存使用情况

	err = pprof.StartCPUProfile(cpuf) // 开始记录 CPU 使用信息
	assert.NoError(t, err)
	defer pprof.StopCPUProfile() // 函数结束后停止记录 CPU 使用情况

	data := make([]int64, 0)
	for i := 0; i < 1e8; i++ {
		data = append(data, int64(i))
	}
	assert.Len(t, data, 1e8)

	heapf, err := os.Create(HEAP_PROFILE_FILENAME)
	assert.NoError(t, err)

	err = pprof.WriteHeapProfile(heapf) // 记录堆内存使用情况
	assert.NoError(t, err)

	time.Sleep(time.Second) // 留有一段记录时间
}
