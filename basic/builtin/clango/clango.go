package clango

/*
#cgo CFLAGS: -O2
#cgo LDFLAGS: -lm -O2

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>     // 引入 Go 语言为调用 C 设置的辅助函数
#include "point.h"      // 可以引用 '.h' 文件, 引入的内容可直接在 Go 代码中调用, 参见 point.h / point.c 文件

char* create_string(const char* ps)
{
    if (!ps)
        return NULL;

    size_t len = strlen(ps) + 1;

    char* pcs = (char*) malloc(len);
    strcpy(pcs, ps);

    return pcs;
}

void show_string(const char* ps)
{
    if (!ps)
        return;

    printf("%s\n", ps);
    fflush(stdout);   // 需要 flush 缓冲区, 否则字符串可能不会显示
}

void free_string(const char* ps)
{
    if (!ps)
        return;

    free((void*)ps);
}
*/
import "C" // C 代码要全部卸载文件最开始的注释中, 紧接着 C 代码注释, 导入 C 代码为 符号 `C`

import (
	"math"
	"unsafe"
)

// # 测试内嵌 C 代码
//
// 调用 `create_string` 函数, 创建 C 字符串
//
// 参数:
//   - `s`: Go 字符串对象
//
// 返回指向 C 字符串的指针
func CreateCString(s string) unsafe.Pointer {
	cs := C.CString(s)         // 将 Go 字符串转为 C 指针
	ptr := C.create_string(cs) // 调用 C 函数, 返回 指针

	return unsafe.Pointer(ptr) // 将指针包装为 Go 对象返回
}

// # 将 C 字符串转为 Go 字符串
//
// 调用 `GoString` C 函数进行转换
//
// 参数:
//   - `ptr`: 指向 C 字符串的指针
//
// 返回 Go 字符串对象
func ConvertCString(ptr unsafe.Pointer) string {
	// Go 语言无法直接使用 C 的字符串, 需要进行转换
	// 传参时, 需要将 unsafe.Pointer 类型参数转为 *C.char 的 C 指针类型
	return C.GoString((*C.char)(ptr))
}

// # 显示 C 字符串
//
// 调用 `show_string` C 函数显示 C 字符串
//
// 参数:
//   - `ptr`: 指向 C 字符串的指针
func ShowCString(ptr unsafe.Pointer) {
	// 传参时, 需要将 unsafe.Pointer 类型参数转为 *C.char 的 C 指针类型
	C.show_string((*C.char)(ptr))
}

// # 释放 C 指针
//
// 调用 `free_string` C 函数释放内存
//
// 参数:
//   - `ptr`: 指向 C 字符串的指针
func FreeCString(ptr unsafe.Pointer) {
	// 传参时, 需要将 unsafe.Pointer 类型参数转为 *C.char 的 C 指针类型
	C.free_string((*C.char)(ptr))
}

// 测试外部 C 代码
// 外部代码是通过 `#include "point.h"` 引入的

// # 包装 C 结构体
//
// 结构体成员变量为一个 C 结构体变量
type Point struct {
	point C.point // 结构体成员为一个 C 结构体
}

// # 创建 `Point` 对象
//
// 通过调用 `create_point` C 函数创建 `point` C 结构体
//
// 参数:
//   - `x`: 结构体成员变量值
//   - `y`: 结构体成员变量值
//
// 返回 `Point` 结构体指针
func CreatePoint(x float64, y float64) *Point {
	// 参数需要转换为 C.double 类型, 返回 point C 结构体变量
	pt := C.create_point((C.double)(x), (C.double)(y))
	return &Point{point: pt}
}

// # 获取 C `point` 结构体的 `x` 成员变量值
//
// 返回 `x` 成员变量值
func (p *Point) GetX() float64 {
	return float64(p.point.x) // Go 中可以直接访问 C 结构体 中的字段
}

// # 获取 C `point` 结构体的 `y` 成员变量值
//
// 返回 `y` 成员变量值
func (p *Point) GetY() float64 {
	return float64(p.point.y)
}

// # 计算两个 `Point` 变量之间的距离
//
// 调用 `distance` C 函数, 计算两个点之间的距离
//
// 参数:
//   - `op`: `Point` 变量地址
//
// 返回当前 `Point` 变量和目标 `Point` 变量距离
func (p *Point) Distance(op *Point) float64 {
	pa := (*C.point)(unsafe.Pointer(&p.point)) // 将 point 结构体的地址转为 C 语言类型指针
	pb := (*C.point)(unsafe.Pointer(&op.point))
	return math.Round(float64(C.distance(pa, pb))*100) / 100 // 计算并保留 2 位小数
}

// 对于库的使用 (静态库或动态库)
// GO 使用 C 的库也很简单, 和引入源文件的方式类似, 只需要在注释中设置 "#cgo CFLAGS:" 和 "#cgo LDFLAGS:" 编译选项, 标明库的位置和链接方式即可
//
// 具体步骤如下:
// 1. 假设具备: test.c 和 test.h 文件
// 2. 编译库文件:
//      2.1. 编译静态库文件
//          $ gcc -c main.c            # 编译 main.c 文件, 生成 main.o 目标文件, '-c' 选项表示只编译, 不进行链接
//          $ ar rcs libtest.a main.o  # 将 main.o 打包为静态库, 库名称必须为 libxxx.a, xxx 为任意名称, 另外, 如果有多个 '.o' 文件, 可以在末尾继续追加, 或者使用 '*.o'
//                                     # 参数: 'r' 替换库中已有的目标文件, 或加入新的目标文件; 'c' 不管库否存在都将创建; 's' 创建文件索引, 能提高速度
//      2.2. 编译动态库文件
//          $ gcc -fPIC -shared test.c -o libtest.so   # 将源文件直接编译为动态库
//          # 或者
//          $ gcc -fPIC -c test.c -o test.o            # 先进行编译, 生成二进制结果, 不链接
//          $ gcc -shared test.o -o libtest.so         # 将编译的二进制文件链接成动态库
// 3. 链接库文件
//          #cgo CFLAGS: -I./include               // 设置头文件位置, 假设 test.h 在 ./include 目录下
//          #cgo LDFLAGS: -L${SRCDIR}/lib -l test  // 设置要链接的静态库, 表示 ./lib/libtest.a 或 ./lib/libtest.so 文件. 注意, 由于链接路径必须为绝对路径, 所以需要用 ${SRCDIR} 变量
//    如果链接动态库, 还要注意: 把 '.so' 文件复制到 '/etc/ld.so.conf.d' 下路径包含的位置, 或者设置 'LD_LIBRARY_PATH' 环境变量
