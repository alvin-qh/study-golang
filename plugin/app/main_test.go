package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试插件加载
func TestLoadPlugin(t *testing.T) {
	// 获取插件路径
	pluginPath := makePluginPath()

	// 备份标准输出对象
	backup := os.Stdout

	// 创建管道
	r, w, err := os.Pipe()
	assert.Nil(t, err)

	defer func() {
		r.Close()
		w.Close()

		os.Stdout = backup
	}()

	// 用管道的 Writer 替换标准输出
	os.Stdout = w

	// 运行插件 1 导入的函数
	runner := loadPlugin(fmt.Sprintf("%s/p1.so", pluginPath))
	runner.Run("Hello")

	// 运行插件 2 导入的函数
	runner = loadPlugin(fmt.Sprintf("%s/p2.so", pluginPath))
	runner.Run("OK")

	// 读取标准输出内容, 确认插件导入函数调用正确
	rd := bufio.NewReader(r)

	line, err := rd.ReadString('\n')
	assert.Nil(t, err)
	assert.Equal(t, "Plugin1 report: Hello\n", line)

	line, err = rd.ReadString('\n')
	assert.Nil(t, err)
	assert.Equal(t, "Plugin2 report: OK\n", line)
}
