//go:build cgo

package external

/*
#cgo CFLAGS: -O2
#cgo LDFLAGS: -lm -O2

#include <unistd.h>     // 引入 Go 语言为调用 C 设置的辅助函数
#include "point.h"      // 引入 C 代码中定义的结构体
*/
import "C"
import "math"

// 为 C 结构体重新映射一个类型名称
type Point C.point

// 创建 `Point` 对象
//
// 通过调用 `create_point` C 函数创建 `point` C 结构体
//
// 返回 `Point` 结构体指针
func CreatePoint(x float64, y float64) *Point {
	// 参数需要转换为 C.double 类型, 返回 point C 结构体变量
	pt := Point(C.create_point((C.double)(x), (C.double)(y)))
	return &pt
}

// 获取 C `point` 结构体的 `x` 成员变量值
//
// 可以看到, 对于 C 语言定义的结构体, 可以在 Go 中直接访问其成员变量
//
// 返回 `x` 成员变量值
func (p *Point) GetX() float64 {
	// Go 中可以直接访问 C 结构体 中的字段
	return float64(Point(*p).x)
}

// 获取 C `point` 结构体的 `y` 成员变量值
//
// 返回 `y` 成员变量值
func (p *Point) GetY() float64 {
	// Go 中可以直接访问 C 结构体 中的字段
	return float64(Point(*p).y)
}

// 计算两个 `Point` 变量之间的距离
//
// 调用 `distance` C 函数, 计算两个点之间的距离
//
// 返回当前 `Point` 变量和目标 `Point` 变量距离
func (p *Point) Distance(op *Point) float64 {
	// 将 point 结构体的地址转为 C 语言类型指针
	pa := (*C.point)(p)
	pb := (*C.point)(op)

	// 计算并保留 2 位小数
	return math.Round(float64(C.distance(pa, pb))*100) / 100
}
