package callstack

import (
	"bytes"
	"runtime"
)

var (
	stackSuffix = []byte("\ngoroutine 1") // 有效调用堆栈的结束标记 (如果存在)
	lineBreak   = []byte("\n")            // 换行标记
)

// 获取调用堆栈
// 调用堆栈的第 1 行为: 协程信息
// 之后内容为:
//
//	调用函数
//	      调用位置
//	调用函数
//	      调用位置
//	...
func CallStack() string {
	stack := make([]byte, 16<<10) // 创建 16kb buffer

	length := runtime.Stack(stack, true) // 获取完整的调用堆栈
	stack = stack[:length]               // 截取调用堆栈的有效长度

	// 查找结束标记是否存在
	// 如果在 panic 之后获取调用堆栈, 则可以去掉 "\ngoroutine" 这一行以及之后的消息
	end := bytes.LastIndex(stack, stackSuffix)
	if end < 0 {
		end = bytes.LastIndex(stack, lineBreak)
	}

	return string(stack[:end])
}
