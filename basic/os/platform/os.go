package platform

import (
	"math/bits"
	"runtime"
	"strings"
	"study/basic/expression"
)

type OSType uint

const (
	Unknown OSType = 0
	Windows OSType = 0x01
	Linux   OSType = 0x02
	Darwin  OSType = 0x04
)

var (
	goosMapping = map[OSType]string{
		Windows: "windows",
		Linux:   "linux",
		Darwin:  "darwin",
	}
)

// 获取当前系统标识
func CurrentOS() OSType {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		return Darwin
	}
	return Unknown
}

// 枚举所有可能的系统类型标识
func (o OSType) each(f func(os OSType) bool) {
	val := o
	for val > 0 {
		t := Unknown
		switch {
		case val&Windows == Windows:
			t = Windows
		case val&Linux == Linux:
			t = Linux
		case val&Darwin == Darwin:
			t = Darwin
		}

		if t == Unknown || !f(t) {
			break
		}
		val &= ^t
	}
}

// 获取标识中第一个系统类型
func (o OSType) First() OSType {
	r := Unknown
	o.each(func(os OSType) bool {
		r = os
		return false
	})
	return r
}

// 判读系统标识是否包含指定的系统类型
func (o OSType) Contains(os OSType) bool {
	return o&os == os
}

// 向系统标识添加系统类型
func (o OSType) Add(os OSType) OSType {
	return o | os
}

// 从标识中删除一个系统类型
func (o OSType) Remove(os OSType) OSType {
	return o &^ os
}

// 获取 OSType 代表的操作系统类型
func (o OSType) List() []OSType {
	r := make([]OSType, 0)
	o.each(func(os OSType) bool {
		r = append(r, os)
		return true
	})
	return r
}

// 将所有的系统标识转为字符串
func (o OSType) String() string {
	switch bits.OnesCount(uint(o)) {
	case 0:
		return ""
	case 1:
		if goos, ok := goosMapping[o]; ok {
			return goos
		}
		return ""
	default:
		r := make([]string, 0)
		o.each(func(os OSType) bool {
			r = append(r, goosMapping[os])
			return true
		})
		return strings.Join(r, " ")
	}
}

// 判断操作系统是否为指定系统
func IsOSMatch(os OSType) bool {
	cur := CurrentOS()

	matched := false
	os.each(func(t OSType) bool {
		if t != cur {
			return true
		}

		matched = true
		return true
	})
	return matched
}

// 判断操作系统是否不为指定系统
func IsOSNotMatch(os OSType) bool {
	return !IsOSMatch(os)
}

// 在指定操作系统上执行回调函数
func RunIfOSMatch(os OSType, fn func()) {
	if IsOSMatch(os) {
		fn()
	}
}

// 在指定操作系统以外的操作系统上执行回调函数
func RunIfOSNot(os OSType, fn func()) {
	if IsOSNotMatch(os) {
		fn()
	}
}

// 根据系统类型选择不同的值
func Choose[T any](os OSType, vMatch T, vMiss T) T {
	return expression.If(IsOSMatch(os), vMatch, vMiss)
}
