module module/app

go 1.21.2

replace module/demo-module => ../module

require (
	gitee.com/go-common-libs/demo-module v1.0.4
	module/demo-module v0.0.0-00010101000000-000000000000
)
