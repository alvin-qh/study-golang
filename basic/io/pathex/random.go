package pathex

import (
	"fmt"
	"math/rand"
	"time"
)

// 随机文件参数选项
type randomFileOpt struct {
	prefix string
	ext    string
}

// 随机文件选项
type RandomFileOption func(*randomFileOpt)

// 设置文件名前缀
func WithPrefix(prefix string) RandomFileOption {
	return func(opt *randomFileOpt) {
		opt.prefix = prefix
	}
}

// 设置文件名后缀
func WithExt(ext string) RandomFileOption {
	return func(opt *randomFileOpt) {
		opt.ext = ext
	}
}

// 返回一个随机文件名
func RandomFileName(opts ...RandomFileOption) string {
	// 定义默认参数
	opt := randomFileOpt{
		prefix: "f-",
		ext:    ".txt",
	}

	// 注入可选参数
	for _, o := range opts {
		o(&opt)
	}

	now := time.Now().UTC()
	return fmt.Sprintf("%s%s%d%s", opt.prefix, now.Format("20060102150405"), rand.Intn(900)+100, opt.ext)
}

// 返回一个随机文件名
func RandomDirName(opts ...RandomFileOption) string {
	// 定义默认参数
	opt := randomFileOpt{
		prefix: "d-",
		ext:    "",
	}

	// 注入可选参数
	for _, o := range opts {
		o(&opt)
	}

	now := time.Now().UTC()
	return fmt.Sprintf("%s%s%d", opt.prefix, now.Format("20060102150405"), rand.Intn(900)+100)
}
