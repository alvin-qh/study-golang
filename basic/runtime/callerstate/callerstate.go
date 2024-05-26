package callerstate

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	ErrNoCallerResolved = errors.New("no caller function resolved")
)

// 表示一个调用帧信息
//
// 一个调用帧即调用栈上一条一条数据, 包含:
//   - 程序计数器 (`pc`)
//   - 调用方所在的文件 (`file`)
//   - 调用方所在的行号 (`line`)
//   - 调用方的函数名称 (`funcName`)
type CallerState struct {
	pc       uintptr
	FileName string
	LineNo   int
	FuncName string
}

// 创建调用帧实例
//
// 通过程序计数器值 (`pc`) 查询对应的调用函数信息, 返回其对应调用帧的完整信息实例
func New(pc uintptr) *CallerState {
	// 根据程序计数器值获取对应函数实例
	fn := runtime.FuncForPC(pc)

	// 获取堆栈函数所在的文件和代码调用位置信息
	file, line := fn.FileLine(pc)

	return &CallerState{
		pc:       pc,
		FileName: filepath.Clean(file),
		LineNo:   line,
		FuncName: fn.Name(),
	}
}

// 将对象转为字符串
func (cs *CallerState) String() string {
	return fmt.Sprintf("%v:%v(%v)", cs.FuncName, cs.FileName, cs.LineNo)
}

// 获取调用方信息
//
// 在当前函数内, 通过 `runtime.Caller(skip)` 函数即可获取调用当前函数的调用方信息,
// 包括: 程序计数器 (`pc`), 调用方所在的文件 (`file`), 调用方所在的行号 (`line`)
//
// `skip` 参数表示要跳过的调用帧数量:
//   - `skip=0` 表示不跳过帧, 即当前函数调用 `runtime.Caller` 函数的调用位置;
//   - `skip=1` 表示跳过一帧, 即调用当前函数的外层函数调用位置;
//   - `skip=n` 以此类推;
func getCallerInfo(skip int) (*CallerState, error) {
	// 获取调用帧信息
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil, ErrNoCallerResolved
	}

	// 根据程序计数器值获取对应函数实例
	fn := runtime.FuncForPC(pc)

	return &CallerState{
		pc:       pc,
		FileName: filepath.Clean(file),
		LineNo:   line,
		FuncName: fn.Name(),
	}, nil
}

// 获取调用该 `Where` 函数的外层函数信息
//
// 该函数内部通过 `runtime.Caller(2)` 函数获取调用当前函数的外层函数信息
func Where() (*CallerState, error) {
	return getCallerInfo(2)
}

// 获取当前代码文件名称
//
// 该函数内部通过 `runtime.Caller(2)` 函数获取调用当前函数的外层函数信息,
// 并返回外层函数源代码所在的文件名
func GetCurrentGoFile() (string, error) {
	cs, err := getCallerInfo(2)
	if err != nil {
		return "", err
	}
	return filepath.Clean(cs.FileName), nil
}

// 列举调用方的调用堆栈
//
// 根据 `layers` 参数, 返回指定层数的调用堆栈信息
func ListStackInfo(layers uint) []*CallerState {
	// 最高获取 `layers` 层堆栈信息
	pcs := make([]uintptr, layers)

	// 这里由于历史原因, skip 从 1 开始, 而 runtime. Caller 中则是从 0 开始
	n := runtime.Callers(2, pcs)

	// 截取正确长度
	pcs = pcs[:n]

	// 创建保存结果的集合
	cs := make([](*CallerState), n)

	// 遍历调用堆栈函数句柄, 获取函数信息
	for i, pc := range pcs {
		cs[i] = New(pc)
	}
	return cs
}
