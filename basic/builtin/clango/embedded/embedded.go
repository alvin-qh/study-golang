//go:build cgo

package embedded

/*
#cgo CFLAGS: -O2
#cgo LDFLAGS: -lm -O2

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>     // 引入 Go 语言为调用 C 设置的辅助函数

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
import "C"
import "unsafe" // C 代码要全部卸载文件最开始的注释中, 紧接着 C 代码注释, 导入 C 代码为 符号 `C`
// 注意, `import "C"` 指令之前不能包含任何除 C 代码之外的注释, 所有的注释都会被认作为内嵌的 C 代码

// 测试内嵌 C 代码
//
// 调用 `create_string` 函数, 通过 `s` 参数表示的 Go 字符串创建 C 字符串, 返回指向 C 字符串的指针
func CreateCString(s string) unsafe.Pointer {
	// 将 Go 字符串转为 C 指针
	cs := C.CString(s)

	// 调用 C 函数, 返回 指针
	ptr := C.create_string(cs)

    // 将指针包装为 Go 对象返回
	return unsafe.Pointer(ptr)
}

// 将 C 字符串转为 Go 字符串
//
// 调用 `GoString` C 函数进行转换, `GoString` 函数是 Go 语言框架为了方便和 C 语言集成提供的工具方法, 位于 `unistd.h` 头文件中
//
// 其中, `ptr` 参数为一个指向 C 字符串的指针, 返回 Go 字符串对象
func ConvertCString(ptr unsafe.Pointer) string {
	// Go 语言无法直接使用 C 的字符串, 需要进行转换
	// 传参时, 需要将 unsafe.Pointer 类型参数转为 *C.char 的 C 指针类型
	return C.GoString((*C.char)(ptr))
}

// 显示 C 字符串
//
// 调用 C 函数 `show_string` 显示 C 字符串
//
// 其中, `ptr` 参数为一个指向 C 字符串的指针
func ShowCString(ptr unsafe.Pointer) {
	// 传参时, 需要将 unsafe.Pointer 类型参数转为 *C.char 的 C 指针类型
	C.show_string((*C.char)(ptr))
}

// 释放 C 指针
//
// 调用 C 函数 `free_string` 释放内存
//
// 在调用 `free_string` 函数时需转为 `*C.char` 指针类型, 表示一个 C 语言 `char` 类型指针
//
// 其中, `ptr` 参数为一个指向 C 字符串的指针
func FreeCString(ptr unsafe.Pointer) {
	// 传参时, 需要将 unsafe.Pointer 类型参数转为 *C.char 的 C 指针类型
	C.free_string((*C.char)(ptr))
}
