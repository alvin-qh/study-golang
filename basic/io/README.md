# Go IO

## 1. Interface

### 1.1. 和读取相关的接口

- `io.Reader`: 基本的读取操作, 从流中读取一系列 bytes
- `io.ReaderAt`: 随机读取操作, 读取流中任意起始位置 bytes
- `io.ByteReader`: 读取 1 个 `byte`
- `io.ByteScanner`: 获取剩余未读取的 `byte`s
- `io.RuneReader`: 读取 1 个 `rune`
- `io.RuneScanner`: 获取剩余未读取的 `rune`s
- `io.WriterTo`: 将内容写入另一个 `io.Writer` 接口对象中
- `io.Seeker`: 随机移动读取指针
- `io.Closer`: 关闭当前 Reader 对象

### 1.2. 和写入相关的接口

- `io.Writer`: 基本的写操作, 在流中顺序写入 `byte`s
- `io.WriterAt`: 随机写操作, 在流的任意位置写入 `byte`s
- `io.StringWriter`: 写入字符串操作
- `io.ReadFrom`: 从另一个 `io.Reader` 对象中读取内容写入当前对象中
- `io.Seeker`: 随机移动读取指针
- `io.Closer`: 关闭当前 `Reader` 对象
