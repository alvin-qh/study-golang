module study-golang/module/app

go 1.21.2

replace study-golang/module/demo-module => ../module2

require (
	gitee.com/go-common-libs/demo-module v1.0.4 // indirect
	study-golang/module/demo-module v1.0.0 // indirect
)
