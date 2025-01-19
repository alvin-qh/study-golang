package fs

import "embed"

var (
	// 将 `./asset` 路径下的内容嵌入为文件系统
	//go:embed asset/*
	STATIC_ASSETS embed.FS
)
