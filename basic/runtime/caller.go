package runtime

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

func ListStackInfo() []*CallerState {
	pcs := make([]uintptr, 100)

	//
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	cs := make([](*CallerState), n)
	for i, pc := range pcs {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)

		cs[i] = &CallerState{
			pc:       pc,
			FileName: file,
			LineNo:   line - 1, // 这里因为历史原因，line 计数从 0 开始，而 runtime.Caller 中则是从 1 开始
			FuncName: fn.Name(),
		}
	}
	return cs
}
