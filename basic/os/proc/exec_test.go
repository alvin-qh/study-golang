package proc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"study/basic/io/pipe"
	"study/basic/os/platform"
	"study/basic/testing/testit"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试获取可执行文件的路径
//
// 获取的位置可以是指定位置或者 $PATH 变量定义的路径
func TestExec_LookPath(t *testing.T) {
	// 定义不同平台下可访问到的系统可执行文件 (位于 $PATH 变量中)
	var fname, target string

	// 获取系统可执行文件的路径
	if platform.IsOSMatch(platform.Windows) {
		fname = "cmd.exe"
		target = "c:\\windows\\system32\\" + fname
	} else {
		fname = "bash"
		target = "/bin/" + fname
	}

	// 获取系统可执行文件的路径
	path, err := exec.LookPath(fname)
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(strings.ToLower(path), target))

	fname = "./temp.exe"
	defer os.Remove("./temp.exe")

	// 创建一个可执行文件并查找它, 这里需要明确的绝对或相对路径, 否则将在 $PATH 变量中查找
	// 在 windows 平台下, `.exe` 扩展名表示该文件为可执行文件
	// 在 Linux 平台下, 0777 权限表示该文件为可执行文件
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0777)
	assert.Nil(t, err)

	f.Close()

	// 查找可执行文件
	path, err = exec.LookPath("./temp.exe")
	assert.Nil(t, err)
	assert.Equal(t, "./temp.exe", path)
}

// 测试执行一个命令
//
// `exec.Command` 为高阶函数, 底层依赖 `os.StartProcess` 函数, 简化了调用流程,
// 可以直接通过命令行字符串执行所需进程
func TestExec_Command(t *testing.T) {
	// 创建一个命令对象, 用于启动一个进程
	cmd := exec.Command("ls", "-l")

	// 创建管道实例接收标准输出
	pi, err := pipe.New()
	assert.Nil(t, err)

	defer pi.Close()

	// 令命令的标准输出指向管道
	cmd.Stdout = pi.Writer()

	// 执行命令
	err = cmd.Run()
	assert.Nil(t, err)

	// 关闭写管道
	pi.CloseWriter()

	// 读取标注输出
	out := bytes.NewBuffer(make([]byte, 0, 1024))
	pi.ReadTo(out)

	fmt.Println(out.String())
}

// 测试执行一个命令, 并传递一个上下文实例
//
// 可以传递一个 `CancelContext` 或者 `TimeoutContext`, 用于在进程执行过程中取消或令进程超时,
// 将会向被取消的进程发送 `os.Kill` 信号, 令进程退出
func TestExec_CommandContext(t *testing.T) {
	testit.SkipTimeOnOS(t, platform.Windows)

	// 创建超时上下文实例
	ctx, cancelFn := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFn()

	// 创建一个命令对象, 用于启动一个进程
	cmd := exec.CommandContext(ctx, "sleep", "10s")

	// 执行命令, 返回错误信息
	err := cmd.Run()
	// 错误表示进程已被杀死
	assert.EqualError(t, err, "signal: killed")
}
