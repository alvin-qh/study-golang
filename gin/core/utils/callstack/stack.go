package callstack

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// 显示源代码文件中指定行数的源代码内容
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// 返回包含调用地址的函数名称
func function(pc uintptr) []byte {
	// 根据调用地址获取调用函数
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())

	// 如果函数名中包含 `/`, 则删除 `/` 之前的内容; 如果函数中包含 `.`, 则删除 `.` 之前的内容; 如果函数中包含 `·`, 则改为 `.`
	// 例如:
	//   runtime/debug.*T·ptrmethod
	// 格式化为
	//   *T.ptrmethod
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.ReplaceAll(name, centerDot, dot)

	return name
}

// 获取格式化后的调用堆栈 (跳过指定的调用帧)
//
// 参数:
//   - `skip` (`int`): 要跳过的调用帧
//
// 返回:
//   - 格式化后的调用堆栈字符串
func CallStack(skip int) []byte {
	buf := new(bytes.Buffer)

	var lines [][]byte
	var lastFile string

	for i := skip; ; i++ {
		// 获取第 i 个调用方
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 输出调用方的文件, 行号和调用地址
		fmt.Fprintf(buf, "\t%s:%d (0x%x)\n", file, line, pc)

		// 如果当前文件内容上次未展示, 则展示文件内容
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t\t%s: %s\n", function(pc), source(lines, line))
	}

	return buf.Bytes()
}
