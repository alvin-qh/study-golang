package platform

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试获取当前操作系统的标识符
func TestOS_CurrentOS(t *testing.T) {
	var goos string
	switch CurrentOS() {
	case Windows:
		goos = "windows"
	case Linux:
		goos = "linux"
	case Darwin:
		goos = "darwin"
	}

	assert.Equal(t, goos, runtime.GOOS)
}

// 获取标识符标识的所有系统
func TestOS_each(t *testing.T) {
	os := Windows | Linux

	l := make([]OSType, 0)
	os.each(func(t OSType) bool {
		l = append(l, t)
		return true
	})

	assert.ElementsMatch(t, []OSType{Windows, Linux}, l)
}

// 测试系统标识中所有的系统类型
func TestOS_List(t *testing.T) {
	os := Windows | Linux
	assert.Equal(t, []OSType{Windows, Linux}, os.List())
}

// 测试系统标识转字符串
func TestOS_String(t *testing.T) {
	os := Windows | Linux
	assert.Equal(t, "windows linux", os.String())
}

// 获取标识符标识的所有系统
func TestOS_First(t *testing.T) {
	os := Windows | Linux
	assert.Equal(t, Windows, os.First())
}

// 测试删除指定类型系统
func TestOS_Remove(t *testing.T) {
	os := Windows | Linux

	os = os.Remove(Linux)
	assert.True(t, os.Contains(Windows))
	assert.False(t, os.Contains(Linux))
}

// 测试添加指定类型系统
func TestOS_Add(t *testing.T) {
	os := Windows | Linux
	assert.False(t, os.Contains(Darwin))

	os = os.Add(Darwin)
	assert.True(t, os.Contains(Windows))
	assert.True(t, os.Contains(Linux))
	assert.True(t, os.Contains(Darwin))
}

// 测试匹配当前操作系统
func TestOS_IsOSMatch(t *testing.T) {
	assert.True(t, IsOSMatch(Windows|Linux|Darwin))

	if runtime.GOOS == "windows" {
		assert.True(t, IsOSMatch(Windows))
		assert.False(t, IsOSMatch(Linux|Darwin))
	} else if runtime.GOOS == "linux" {
		assert.True(t, IsOSMatch(Linux))
		assert.False(t, IsOSMatch(Windows|Darwin))
	} else if runtime.GOOS == "darwin" {
		assert.True(t, IsOSMatch(Darwin))
		assert.False(t, IsOSMatch(Windows|Linux))
	} else {
		t.Error("unknown os")
	}
}

// 测试匹配不是当前操作系统
func TestOS_IsOSNotMatch(t *testing.T) {
	assert.False(t, IsOSNotMatch(Windows|Linux|Darwin))

	if runtime.GOOS == "windows" {
		assert.True(t, IsOSNotMatch(Linux|Darwin))
		assert.False(t, IsOSNotMatch(Windows))
	} else if runtime.GOOS == "linux" {
		assert.True(t, IsOSNotMatch(Windows|Darwin))
		assert.False(t, IsOSNotMatch(Linux))
	} else if runtime.GOOS == "darwin" {
		assert.True(t, IsOSNotMatch(Windows|Linux))
		assert.False(t, IsOSNotMatch(Darwin))
	} else {
		t.Error("unknown os")
	}
}

// 测试当匹配操作系统时进行回调
func TestOS_RunIfOSMatch(t *testing.T) {
	matched := false

	// 测试匹配当前系统标识
	if runtime.GOOS == "windows" {
		RunIfOSMatch(Windows, func() {
			matched = true
		})
	} else if runtime.GOOS == "linux" {
		RunIfOSMatch(Linux, func() {
			matched = true
		})
	} else if runtime.GOOS == "darwin" {
		RunIfOSMatch(Darwin, func() {
			matched = true
		})
	} else {
		t.Error("unknown os")
	}
	assert.True(t, matched)

	matched = false

	// 测试在多个标识中匹配当前系统标识
	if runtime.GOOS == "windows" {
		RunIfOSMatch(Windows|Darwin, func() {
			matched = true
		})
	} else if runtime.GOOS == "linux" {
		RunIfOSMatch(Linux|Windows, func() {
			matched = true
		})
	} else if runtime.GOOS == "darwin" {
		RunIfOSMatch(Darwin|Linux, func() {
			matched = true
		})
	} else {
		t.Error("unknown os")
	}
	assert.True(t, matched)

	matched = false

	// 匹配不到当前系统标识
	if runtime.GOOS == "windows" {
		RunIfOSMatch(Linux|Darwin, func() {
			matched = true
		})
	} else if runtime.GOOS == "linux" {
		RunIfOSMatch(Darwin|Windows, func() {
			matched = true
		})
	} else if runtime.GOOS == "darwin" {
		RunIfOSMatch(Windows|Linux, func() {
			matched = true
		})
	} else {
		t.Error("unknown os")
	}
	assert.False(t, matched)
}

// 测试根据操作系统标识获取不同的值
func TestOS_Choose(t *testing.T) {
	if runtime.GOOS == "windows" {
		assert.Equal(t, "windows", Choose(Windows, "windows", ""))
	} else if runtime.GOOS == "linux" {
		assert.Equal(t, "linux", Choose(Linux, "linux", ""))
	} else if runtime.GOOS == "darwin" {
		assert.Equal(t, "darwin", Choose(Darwin, "darwin", ""))
	} else {
		t.Error("unknown os")
	}
}
