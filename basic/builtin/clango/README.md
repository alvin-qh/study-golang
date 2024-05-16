# C 链接库使用

## 1. 使用方法

GO 使用 C 的库也很简单, 和引入源文件的方式类似, 只需要在注释中设置 `#cgo CFLAGS:` 和 `#cgo LDFLAGS:` 编译选项, 标明库的位置和链接方式即可

## 2. 具体步骤

### 2.1. 创建测试文件

创建 `test.c` 和 `test.h` 文件

### 2.2. 编译库文件

#### 2.2.1. 编译静态库文件

```bash
# 编译 main.c 文件, 生成 main.o 目标文件, '-c' 选项表示只编译, 不进行链接
gcc -c main.c

# 将 main.o 打包为静态库, 库名称必须为 libxxx.a, xxx 为任意名称
# 另外, 如果有多个 '.o' 文件, 可以在末尾继续追加, 或者使用 '*.o'
ar rcs libtest.a main.o
# 参数:
#   - 'r' 替换库中已有的目标文件, 或加入新的目标文件;
#   - 'c' 不管库否存在都将创建;
#   - 's' 创建文件索引, 能提高速度
```

#### 2.2.2. 编译动态库文件

```bash
# 将源文件直接编译为动态库
gcc -fPIC -shared test.c -o libtest.so
```

或者

```bash
# 先进行编译, 生成二进制结果, 不链接
gcc -fPIC -c test.c -o test.o

# 将编译的二进制文件链接成动态库
gcc -shared test.o -o libtest.so
```

#### 2.2.3. 链接库文件

假设 `.h` 头文件在 `./include` 目录下

静态链接库的路径必须为绝对路径, 可以使用 `${SRCDIR}` 表示当前源码的绝对路径

链接动态库需要把 `.so` 文件复制到 `/etc/ld.so.conf.d` 下路径包含的位置, 或者设置 `LD_LIBRARY_PATH` 环境变量

下面例子中表示链接到 `./lib/libtest.a` 或 `./lib/libtest.so` 库

```go
/*
#cgo CFLAGS: -I./include  // 设置头文件位置
#cgo LDFLAGS: -L${SRCDIR}/lib -l test  // 设置要链接的静态库
*/
```
