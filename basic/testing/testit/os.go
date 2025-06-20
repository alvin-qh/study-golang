package testit

import (
	"basic/os/platform"
	"flag"
	"fmt"
	"testing"
)

// 在指定的操作系统上跳过测试
func SkipTimeOnOS(t *testing.T, os platform.OSType) {
	if platform.IsOSMatch(os) {
		t.Skipf("Skipping test on %s", os)
	}
}

// 在指定的操作系统上执行测试
func RunIf(t *testing.T, os platform.OSType, f func(t *testing.T)) bool {
	if platform.IsOSNotMatch(os) {
		return false
	}

	t.Run(fmt.Sprintf("run on %s", os), f)
	return false
}

// 判断当前环境是否测试环境
func IsTesting() bool {
	return flag.Lookup("test.v") != nil
}
