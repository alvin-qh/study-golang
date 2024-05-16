package proc

import (
	"os"
	"os/user"
	"strconv"
	"study/basic/os/platform"
	"study/basic/testing/testit"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 获取启动当前进程的真实用户 UID
//
// 启动进程的用户 ID 应该和当前登陆系统的用户 ID 相同
func TestOS_Getuid(t *testing.T) {
	uid := os.Getuid()

	user, err := user.Current()
	assert.Nil(t, err)

	assert.Equal(t, user.Uid, strconv.Itoa(uid))
}

// 获取启动当前进程的有效用户 UID
//
// 有效用户一般情况下和真实用户一致, 除非用户通过 `sudo` 类命令将用户临时改为 `root` 用户
func TestOS_Geteuid(t *testing.T) {
	testit.SkipTimeOnOS(t, platform.Windows)

	uid := os.Geteuid()
	assert.Equal(t, os.Getuid(), uid)
}

// 获取启动当前进程的真实用户所在组的 GID
//
// 启动进程的用户组 ID 应该和当前登陆系统的用户组 ID 相同
func TestOS_Getgid(t *testing.T) {
	gid := os.Getgid()

	user, err := user.Current()
	assert.Nil(t, err)

	assert.Equal(t, user.Gid, strconv.Itoa(gid))
}

// 获取启动当前进程的有效用户组 GID
//
// 有效用户组一般情况下和真实用户组一致, 除非用户通过 `sudo` 类命令将用户临时改为 `root` 用户
func TestOS_Getegid(t *testing.T) {
	testit.SkipTimeOnOS(t, platform.Windows)

	gid := os.Getegid()
	assert.Equal(t, os.Getgid(), gid)
}
