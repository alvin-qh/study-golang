// 针对 Linux 平台编译
package cross

import (
	"os"
	"syscall"
	"time"
)

// Linux 平台下, 从 `os.FileInfo` 获取文件的 `uid` 和 `gid`
func FileOwner(fi os.FileInfo) (uint32, uint32, bool) {
	stat_t, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, false
	}

	return stat_t.Uid, stat_t.Gid, true
}

// Linux 平台下获取文件的最后访问时间
//
// Go 语言框架并没有直接获取文件访问时间的 API, 需要通过各系统平台的系统调用获取
func FileAtime(fi os.FileInfo) (time.Time, bool) {
	stat_t, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return time.UnixMilli(0), false
	}

	sec, nsec := stat_t.Atim.Unix()
	return time.Unix(sec, nsec), true
}
