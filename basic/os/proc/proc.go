package proc

import (
	"bytes"
	"fmt"
	"os"
)

// 获取当前进程命令行 (For Linux)
//
// Go 语言并没有提供获取进程命令行的方法, 在 Linux 系统上可以通过访问 `/proc/pid/cmdline` 文件获取当前进程命令行
// 对于 Windows 系统, 则只能通过调用 Win32 API 获取进程命令行
func GetCmdLinesByPid(pid int) ([]string, error) {
	// 获取当前进程命令行
	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)
	// 将命令行分解, 并存入切片
	for _, line := range bytes.Split(cmdline, []byte("\x00")) {
		lines = append(lines, string(bytes.TrimSpace(line)))
	}
	return lines, nil
}
