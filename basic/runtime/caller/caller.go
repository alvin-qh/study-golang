package caller

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	ErrNoCallerResolved = errors.New("no caller function resolved")
)

type CallerState struct {
	pc       uintptr
	FileName string
	LineNo   int
	FuncName string
}

// 转为字符串
func (cs *CallerState) String() string {
	return fmt.Sprintf("%v:%v(%v)", cs.FuncName, cs.FileName, cs.LineNo)
}

//
func getCallerInfo(skip int) (*CallerState, error) {
	// 获取 caller，即调用方信息
	// 参数 skip 表示跳过多少层的调用，0 表示当前函数，1 表示上一层函数，以此类推
	// pc 表示调用方函数的句柄；file 表示调用方所在的文件；line 表示调用方的行号；ok 表示是否成功
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil, ErrNoCallerResolved
	}

	// 获取调用方函数的信息对象
	fn := runtime.FuncForPC(pc)

	return &CallerState{
		pc:       pc,
		FileName: file,
		LineNo:   line,
		FuncName: fn.Name(), // 获取被调用函数
	}, nil
}

// 获取调用该函数额调用方信息
func Where() (*CallerState, error) {
	return getCallerInfo(2)
}

// 获取当前代码文件名称
func GetCurrentGoFile() (string, error) {
	cs, err := getCallerInfo(2)
	if err != nil {
		return "", err
	}
	return cs.FileName, nil
}

// 列举调用方的调用堆栈
func ListStackInfo() []*CallerState {
	pcs := make([]uintptr, 100) // 最高获取 100 层堆栈信息

	n := runtime.Callers(2, pcs) // 这里由于历史原因，skip 从 1 开始，而 runtime。Caller 中则是从 0 开始
	pcs = pcs[:n]                // 截取正确长度

	cs := make([](*CallerState), n) // 保持结果的 slice

	// 遍历调用堆栈函数句柄，获取函数信息
	for i, pc := range pcs {
		fn := runtime.FuncForPC(pc)   // 从函数句柄获取函数信息
		file, line := fn.FileLine(pc) // 获取堆栈函数所在的文件和代码调用位置信息

		// 填充结果 slice
		cs[i] = &CallerState{
			pc:       pc,
			FileName: file,
			LineNo:   line - 1, // line 和实际行号差 1
			FuncName: fn.Name(),
		}
	}
	return cs
}
