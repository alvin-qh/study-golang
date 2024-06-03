package proc

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"study/basic/io/pipe"
	"study/basic/os/platform"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试获取当前进程的执行文件名
func TestOS_Executable(t *testing.T) {
	n, err := os.Executable()
	assert.Nil(t, err)

	if platform.IsOSMatch(platform.Windows) {
		assert.True(t, strings.HasSuffix(n, ".exe"))
	} else {
		assert.True(t, strings.HasSuffix(n, "/proc.test"))
	}
}

// 测试获取当前进程的进程 ID
func TestOS_Getpid(t *testing.T) {
	// 获取当前进程的进程 ID
	pid := os.Getpid()
	assert.Greater(t, pid, 0)

	if platform.IsOSNotMatch(platform.Windows) {
		cmdlines, err := GetCmdLinesByPid(pid)
		assert.Nil(t, err)
		assert.Greater(t, len(cmdlines), 0)

		// 确认进程命令行的第一个参数为当前进程可执行文件路径
		assert.True(t, strings.HasSuffix(cmdlines[0], "/proc.test"))
		// assert.Contains(t, lines, "-test.timeout=30s")
		// assert.Contains(t, lines, "-test.count=1")
		// assert.Contains(t, lines, "-test.v=true")
	}
}

// 测试获取父进程的进程 ID
func TestOS_Getppid(t *testing.T) {
	// 获取当前进程父进程的 ID
	pid := os.Getppid()
	assert.Greater(t, pid, 0)

	if platform.IsOSNotMatch(platform.Windows) {
		// 获取父进程的命令行
		cmdlines, err := GetCmdLinesByPid(pid)
		assert.Nil(t, err)
		assert.Greater(t, len(cmdlines), 0)

		// 确认进程命令行的第一个参数为当前进程可执行文件路径
		assert.True(t, strings.HasSuffix(cmdlines[0], "go"))
		assert.Equal(t, "test", cmdlines[1])
	}
}

// 测试启动进程
//
// `os.StartProcess` 是进程处理的低阶 API, 用于一些特殊进程的处理, 一般情况下, 应该使用
// `exec.Command` 这样的高阶 API
func TestOS_StartProcess(t *testing.T) {
	// 定义一个管道, 用于接收 stdout 内容
	pi, err := pipe.New()
	assert.Nil(t, err)

	defer pi.Close()

	// 定义进程可执行文件路径
	cmd := platform.Choose(
		platform.Windows,
		"C:\\Windows\\system32\\cmd.exe",
		"ls",
	)
	args := platform.Choose(
		platform.Windows,
		[]string{"C:\\Windows\\system32\\cmd.exe", "/C", "dir"},
		[]string{"ls"},
	)
	attr := &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			pi.Writer(),
			os.Stderr,
		},
	}

	// 启动进程
	// 注意, 进程参数的第一项必须是进程对应的文件名称
	p, err := os.StartProcess(cmd, args, attr)
	assert.Nil(t, err)

	// 等待进程执行完毕
	fi, err := p.Wait()
	assert.Nil(t, err)

	// 进程结束, 关闭管道写通道
	pi.CloseWriter()

	// 获取进程标准输出内容
	out := bytes.NewBuffer(make([]byte, 0, 1024))
	err = pi.ReadTo(out)
	assert.Nil(t, err)

	// 确认进程正确结束
	assert.True(t, fi.Exited())
	assert.True(t, fi.Success())
	assert.Equal(t, p.Pid, fi.Pid())
	assert.Equal(t, 0, fi.ExitCode())
	// 进程执行占用的系统 CPU 时间
	assert.Equal(t, int64(0), fi.SystemTime().Nanoseconds())
	// 系统执行占用的用户 CPU 时间
	assert.GreaterOrEqual(t, fi.UserTime().Nanoseconds(), int64(0))

	fmt.Println(out.String())
}

// 测试杀死启动的进程
//
// 通过 `os.Process` 实例的 `Kill` 方法可以杀死其代表的进程,
// 被杀死的进程具备特殊的退出码
func TestOS_Process_Kill(t *testing.T) {
	// 定义进程可执行文件路径
	cmd := platform.Choose(platform.Windows, "C:\\Windows\\system32\\ping.exe", "/usr/bin/sleep")
	args := platform.Choose(
		platform.Windows,
		[]string{"C:\\Windows\\system32\\ping.exe", "-n", "10", "127.0.0.1"},
		[]string{"sleep", "10s"},
	)
	attr := &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}

	// 启动进程, 创建进程实例, 该进程会执行 10 秒后退出
	p, err := os.StartProcess(cmd, args, attr)
	assert.Nil(t, err)

	// 异步函数, 在 500 毫秒后杀死进程
	go func() {
		time.Sleep(500 * time.Millisecond)
		p.Kill()
	}()

	// 等待进程执行结束
	fi, err := p.Wait()
	assert.Nil(t, err)

	assert.Equal(t, p.Pid, fi.Pid())
	// 进程未能正确退出
	if platform.IsOSMatch(platform.Windows) {
		assert.True(t, fi.Exited())
	} else {
		assert.False(t, fi.Exited())
	}
	assert.False(t, fi.Success())
	// 确认进程的返回码
	if platform.IsOSMatch(platform.Windows) {
		assert.Equal(t, 1, fi.ExitCode())
	} else {
		assert.Equal(t, -1, fi.ExitCode())
	}
}

// 测试向进程发送信号
//
// 通过 `os.Process` 实例的 `Signal` 方法可以向进程发送信号, Go 语言定义了两种进程信号:
//   - `os.Interrupt`, 中断进程执行信号
//   - `os.Kill`, 杀死进程信号, 相当于执行 `os.Process` 实例的 `Kill` 方法
func TestOS_Process_Signal(t *testing.T) {
	// 定义进程可执行文件路径
	cmd := platform.Choose(platform.Windows, "C:\\Windows\\system32\\ping.exe", "/usr/bin/sleep")
	args := platform.Choose(
		platform.Windows,
		[]string{"C:\\Windows\\system32\\ping.exe", "-n", "10", "127.0.0.1"},
		[]string{"sleep", "10s"},
	)
	attr := &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}

	// 启动进程, 创建进程实例, 该进程会执行 10 秒后退出
	p, err := os.StartProcess(cmd, args, attr)
	assert.Nil(t, err)

	// 异步函数, 在 500 毫秒后向进程发送中断信号
	go func() {
		time.Sleep(500 * time.Millisecond)
		if platform.IsOSMatch(platform.Windows) {
			p.Signal(os.Kill)
		} else {
			p.Signal(os.Interrupt)
		}
	}()

	// 等待进程执行结束
	fi, err := p.Wait()
	assert.Nil(t, err)

	assert.Equal(t, p.Pid, fi.Pid())
	// 进程未能正确退出
	if platform.IsOSMatch(platform.Windows) {
		assert.True(t, fi.Exited())
	} else {
		assert.False(t, fi.Exited())
	}
	assert.False(t, fi.Success())
	// 确认进程的返回码
	if platform.IsOSMatch(platform.Windows) {
		assert.Equal(t, 1, fi.ExitCode())
	} else {
		assert.Equal(t, -1, fi.ExitCode())
	}
}
