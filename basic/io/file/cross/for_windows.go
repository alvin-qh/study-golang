// 针对 Windows 平台编译
package cross

import (
	"os"
	"syscall"
	"time"
)

// Windows 平台下, 从 `os.FileInfo` 获取文件的 `uid` 和 `gid`
//
// Windows 平台下文件没有 `uid` 和 `gid` 的概念
func FileOwner(_ os.FileInfo) (uint32, uint32, bool) {
	return 0, 0, false
}

// Windows 平台下获取文件的最后访问时间
//
// Go 语言框架并没有直接获取文件访问时间的 API, 需要通过各系统平台的系统调用获取
func FileAtime(fi os.FileInfo) (time.Time, bool) {
	atime := fi.Sys().(*syscall.Win32FileAttributeData).LastAccessTime
	return time.Unix(0, atime.Nanoseconds()), true
}
