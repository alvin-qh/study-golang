# Go 语言学习项目完整性分析报告

> 分析日期：2026-02-26
> 项目路径：`/home/alvin/Workspace/Study/study-golang/basic`
> Go 版本：1.26.0
> Oracle 审查状态：✅ 已完成

---

## 一、项目概况

### 1.1 统计数据

| 指标         | 数值                         |
| ------------ | ---------------------------- |
| 总 Go 文件数 | ~146 个（非测试） + 测试文件 |
| 测试文件数   | 82 个                        |
| 代码总行数   | ~17,275 行                   |
| 测试覆盖率   | 大部分模块有测试             |
| 测试失败数   | 3 个测试失败                 |

### 1.2 目录结构

```
basic/
├── builtin/          # Go 内置类型和函数
├── concurrency/      # 并发编程（含详细 GMP 文档）
├── container/        # 容器数据结构
├── expression/       # 表达式工具
├── io/              # 输入输出
├── logs/            # 日志实现
├── mod/             # Go Module
├── net/             # 网络编程
├── oop/             # 面向对象编程
├── os/              # 操作系统交互
├── runtime/         # 运行时
├── testing/         # 测试工具
└── time/            # 时间处理
```

### 1.3 依赖项

```
github.com/stretchr/testify v1.11.1  # 测试断言库
github.com/joho/godotenv v1.5.1      # 环境变量
golang.org/x/exp v0.0.0-20260212183809-81e46e3db34a
golang.org/x/sync v0.19.0
github.com/google/uuid v1.6.0
```

---

## 二、现有模块分析

### 2.1 已完成且质量良好的模块 ✅

| 模块                  | 内容                                           | 评价                              | 完成度 |
| --------------------- | ---------------------------------------------- | --------------------------------- | ------ |
| **concurrency/**      | goroutine、channel、context、sync 原语、atomic | 非常完善，含详细的 GMP 调度器文档 | 95%    |
| **builtin/slices/**   | 切片定义、操作、转换                           | 完整的测试覆盖                    | 90%    |
| **builtin/strings/**  | 字符串操作、转换、unicode                      | 内容丰富                          | 90%    |
| **builtin/reflects/** | 反射类型、值、结构体映射                       | 深度适中                          | 85%    |
| **container/**        | Set、List、OrderedMap                          | 实用性强                          | 85%    |
| **net/tcp/**          | TCP 客户端/服务器                              | 完整实现                          | 90%    |
| **net/udp/**          | UDP 客户端/服务器                              | 完整实现                          | 90%    |
| **net/smtp/**         | SMTP 邮件发送                                  | 可运行示例                        | 85%    |
| **io/file/**          | 文件操作（含跨平台）                           | 考虑周全                          | 90%    |
| **io/archive/**       | tar/zip/gzip 归档                              | 标准库学习                        | 85%    |
| **runtime/gc/**       | GC 内存管理                                    | 有实际演示                        | 80%    |
| **time/**             | 时间格式化、解析、时区                         | 测试完善                          | 80%    |

### 2.2 需要补充实现的模块 ⚠️

| 模块          | 当前状态             | 问题             | 优先级 |
| ------------- | -------------------- | ---------------- | ------ |
| **net/http/** | 仅有测试文件         | 缺少独立实现示例 | P1     |
| **net/ws/**   | 空文件 `package ws`  | 完全未实现       | P4     |
| **net/rpc/**  | 空文件 `package rpc` | 完全未实现       | P4     |
| **mod/**      | README 仅标题        | 内容缺失         | P1     |

### 2.3 测试失败的模块 ❌

| 模块                                      | 测试                            | 错误原因                                  | 修复方案                       |
| ----------------------------------------- | ------------------------------- | ----------------------------------------- | ------------------------------ |
| `builtin/strings/strings_test.go`         | `TestStrings_HasSuffix`         | 测试断言失败                              | 检查第 112 行断言逻辑          |
| `builtin/reflects/values/value_test.go`   | `TestReflect_GetValue`          | 第 85 行 `tv.CanSet()` 对值类型返回 false | 删除或改为 `assert.False`      |
| `runtime/callerstate/callerstate_test.go` | `TestCallerState_ListStackInfo` | 硬编码行号 68，实际为 69                  | 使用范围检查或移除精确行号断言 |

**测试修复代码示例：**

```go
// builtin/reflects/values/value_test.go:85
// 错误：assert.True(t, tv.CanSet()) 对值类型永远返回 false
// 修复：删除该行或改为
assert.False(t, tv.CanSet())

// 或使用指针获取可设置的值
tv := reflect.ValueOf(&obj).Elem()
assert.True(t, tv.CanSet())

// runtime/callerstate/callerstate_test.go:70
// 错误：硬编码行号断言
// 修复：改为范围检查
assert.GreaterOrEqual(t, actual, 65)
assert.LessOrEqual(t, actual, 75)
```

---

## 三、缺失内容分析

### 3.1 高优先级 - 核心概念缺失

#### 3.1.1 基础语法补充

| 主题             | 说明                         | 建议位置             | 预计工作量 |
| ---------------- | ---------------------------- | -------------------- | ---------- |
| **控制流语句**   | for、if、switch、select 详解 | `builtin/control/`   | 4-6 小时   |
| **类型系统深入** | 类型别名、类型定义、底层类型 | `builtin/types/`     | 3-4 小时   |
| **常量与枚举**   | iota、常量组、无类型常量     | `builtin/constants/` | 2-3 小时   |
| **指针详解**     | 指针与值语义、逃逸分析       | `builtin/pointers/`  | 3-4 小时   |

#### 3.1.2 标准库学习

| 标准库            | 内容                        | 建议位置          | 预计工作量 |
| ----------------- | --------------------------- | ----------------- | ---------- |
| **encoding/json** | JSON 编解码、标签、流式处理 | `io/serialize/`   | 2-3 小时   |
| **encoding/xml**  | XML 处理                    | `io/serialize/`   | 2 小时     |
| **encoding/csv**  | CSV 文件处理                | `io/csv/`         | 2 小时     |
| **crypto/**       | 哈希、加密基础              | `crypto/`         | 3-4 小时   |
| **regexp**        | 正则表达式                  | `builtin/regexp/` | 3-4 小时   |
| **sort**          | 排序算法、自定义排序        | `builtin/sort/`   | 2-3 小时   |
| **math**          | 数学运算、随机数            | `builtin/math/`   | 2-3 小时   |

#### 3.1.3 并发进阶主题

| 主题               | 说明                                  | 建议位置                       | 预计工作量 |
| ------------------ | ------------------------------------- | ------------------------------ | ---------- |
| **sync.Pool 详解** | 对象池原理与实践                      | `concurrency/sync/pools/`      | 已有部分   |
| **并发模式**       | worker pool、fan-out/fan-in、pipeline | `concurrency/patterns/`        | 6-8 小时   |
| **select 原理**    | 多路复用机制                          | `concurrency/goroutine/chans/` | 2-3 小时   |

### 3.2 中优先级 - 实用技能补充

#### 3.2.1 工程化能力

| 主题           | 说明                        | 建议位置   | 预计工作量 |
| -------------- | --------------------------- | ---------- | ---------- |
| **项目结构**   | 标准项目布局、组织方式      | `mod/`     | P1         |
| **依赖管理**   | go.mod、go.sum、版本选择    | `mod/`     | 2-3 小时   |
| **构建与编译** | go build、ldflags、交叉编译 | `build/`   | 2-3 小时   |
| **代码生成**   | go generate、stringer       | `codegen/` | 2-3 小时   |

#### 3.2.2 调试与测试

| 主题           | 说明               | 建议位置             | 预计工作量 |
| -------------- | ------------------ | -------------------- | ---------- |
| **Benchmark**  | 基准测试、性能分析 | `testing/benchmark/` | 2-3 小时   |
| **pprof**      | CPU/内存性能分析   | `runtime/profile/`   | 已有部分   |
| **覆盖率分析** | go test -cover     | `testing/coverage/`  | 1-2 小时   |
| **Fuzzing**    | 模糊测试           | `testing/fuzz/`      | 2-3 小时   |

#### 3.2.3 数据库操作

| 主题             | 说明             | 建议位置          | 预计工作量 |
| ---------------- | ---------------- | ----------------- | ---------- |
| **database/sql** | SQL 数据库操作   | `database/sql/`   | 4-6 小时   |
| **Redis**        | Redis 客户端使用 | `database/redis/` | 3-4 小时   |

#### 3.2.4 网络进阶

| 主题            | 说明                      | 建议位置    | 预计工作量 |
| --------------- | ------------------------- | ----------- | ---------- |
| **HTTP 服务端** | 路由、中间件、RESTful API | `net/http/` | 4-6 小时   |
| **HTTP 客户端** | 请求、响应、超时控制      | `net/http/` | 2-3 小时   |
| **WebSocket**   | 实时通信                  | `net/ws/`   | 4-6 小时   |
| **gRPC**        | RPC 框架实践              | `net/rpc/`  | 4-6 小时   |
| **TLS/HTTPS**   | 安全通信                  | `net/tls/`  | 3-4 小时   |

### 3.3 低优先级 - 高级主题补充

| 主题          | 说明                     | 建议位置            | 预计工作量 |
| ------------- | ------------------------ | ------------------- | ---------- |
| **泛型深入**  | 类型参数、约束、类型推断 | `builtin/generics/` | 4-6 小时   |
| **插件系统**  | plugin 包使用            | `plugin/`           | 3-4 小时   |
| **unsafe 包** | 不安全操作、内存布局     | `unsafe/`           | 3-4 小时   |
| **CGO 进阶**  | C 库调用、回调           | `builtin/clango/`   | 已有基础   |
| **汇编基础**  | Go 汇编、性能优化        | `asm/`              | 6-8 小时   |

---

## 四、改进建议

### 4.1 测试修复（P0 - 最高优先级）

```go
// 1. builtin/strings/strings_test.go:112
// 检查 HasSuffix 测试断言逻辑

// 2. builtin/reflects/values/value_test.go:85
// 修复方案 A：删除错误的断言
// assert.True(t, tv.CanSet())  // 删除此行

// 修复方案 B：改为正确的断言
assert.False(t, tv.CanSet())

// 修复方案 C：使用指针获取可设置的值
tv := reflect.ValueOf(&obj).Elem()
assert.True(t, tv.CanSet())

// 3. runtime/callerstate/callerstate_test.go:70
// 修复：使用范围检查代替硬编码行号
assert.GreaterOrEqual(t, actual, 65)
assert.LessOrEqual(t, actual, 75)
```

### 4.2 模块补充建议

#### net/http/ 实现建议

```go
// 建议添加:
// - server.go: HTTP 服务器基础示例
// - client.go: HTTP 客户端请求示例
// - middleware.go: 中间件模式示例
// - handler.go: 处理器函数示例
// - restful.go: RESTful API 示例
```

#### mod/ 内容建议

```markdown
# Go Module

## 1. 模块初始化
go mod init

## 2. 依赖管理
go get, go mod tidy, go mod verify

## 3. 版本控制
语义化版本、伪版本、indirect 依赖

## 4. 私有模块
GOPRIVATE, GONOSUMDB, GOPROXY

## 5. 工作区
go.work, 多模块管理

## 6. 常用命令
go mod graph, go mod why, go mod download
```

### 4.3 工程化改进建议

#### 添加 Makefile

```makefile
.PHONY: test lint fmt build

test:
 go test ./... -v -race

lint:
 golangci-lint run

fmt:
 go fmt ./...
 goimports -w .

build:
 go build -o bin/basic .

coverage:
 go test ./... -coverprofile=coverage.out
 go tool cover -html=coverage.out
```

#### 添加 CI 配置 (.github/workflows/ci.yml)

```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - run: go test ./... -v -race
      - run: go vet ./...
```

### 4.4 新增模块建议

```
basic/
├── builtin/
│   ├── constants/     # 新增: 常量与 iota
│   ├── control/       # 新增: 控制流语句
│   ├── pointers/      # 新增: 指针深入
│   ├── regexp/        # 新增: 正则表达式
│   ├── sort/          # 新增: 排序
│   └── generics/      # 新增: 泛型深入
├── concurrency/
│   └── patterns/      # 新增: 并发模式
├── crypto/            # 新增: 加密与哈希
├── database/          # 新增: 数据库操作
│   ├── sql/
│   └── redis/
├── build/             # 新增: 构建与编译
└── codegen/           # 新增: 代码生成
```

---

## 五、学习路径建议

### 阶段一：基础语法（1-2 周）

1. `builtin/arrays/` - 数组
2. `builtin/slices/` - 切片
3. `builtin/maps/` - Map
4. `builtin/strings/` - 字符串
5. `builtin/functions/` - 函数
6. `builtin/errors/` - 错误处理
7. `builtin/defers/` - defer
8. `builtin/constants/` - 常量（待补充）

### 阶段二：类型系统（1 周）

1. `builtin/types/` - 类型系统
2. `builtin/reflects/` - 反射
3. `oop/` - 面向对象
4. `builtin/generics/` - 泛型（待补充）

### 阶段三：并发编程（2-3 周）

1. `concurrency/goroutine/` - Goroutine
2. `concurrency/goroutine/chans/` - Channel
3. `concurrency/sync/` - 同步原语
4. `concurrency/atomic/` - 原子操作
5. `concurrency/goroutine/context/` - Context
6. `concurrency/patterns/` - 并发模式（待补充）

### 阶段四：I/O 与网络（2 周）

1. `io/` - 文件与 I/O
2. `io/serialize/` - 序列化
3. `net/tcp/` - TCP
4. `net/udp/` - UDP
5. `net/http/` - HTTP（待补充）
6. `io/archive/` - 归档

### 阶段五：运行时与工具（1 周）

1. `runtime/` - GC、调用栈
2. `testing/` - 测试
3. `logs/` - 日志
4. `mod/` - 模块管理（待补充）

### 阶段六：实践项目（2 周）

1. `database/` - 数据库操作（待补充）
2. `crypto/` - 加密（待补充）
3. 综合项目练习

---

## 六、优先级行动清单

| 优先级 | 任务                         | 预计工作量  | 备注         |
| ------ | ---------------------------- | ----------- | ------------ |
| **P0** | 修复 3 个测试失败            | 1-2 小时    | 阻断性问题   |
| **P1** | 补充 `net/http/` 示例        | 4-6 小时    | 最常用标准库 |
| **P1** | 完善 `mod/` 文档             | 2-3 小时    | 工程化基础   |
| **P2** | 添加 `database/sql/` 示例    | 4-6 小时    | 实际项目刚需 |
| **P2** | 添加 `concurrency/patterns/` | 6-8 小时    | 并发最佳实践 |
| **P2** | 添加 `builtin/constants/`    | 2-3 小时    | 基础语法补充 |
| **P3** | 补充 `crypto/` 加密哈希      | 3-4 小时    | 安全相关     |
| **P3** | 添加 benchmark 示例          | 2-3 小时    | 性能意识培养 |
| **P4** | `net/ws/`、`net/rpc/`        | 各 4-6 小时 | 特定场景需求 |

---

## 七、总结

### 优点

- ✅ 并发模块内容详实，GMP 调度器文档（`concurrency/README.md`）是一大亮点
- ✅ 测试覆盖率高，大部分代码都有对应测试
- ✅ 项目结构清晰，模块划分合理
- ✅ time 模块测试完善，涵盖格式化、解析、时区等核心功能

### 待改进

- ⚠️ 存在 3 个测试失败需要修复（P0）
- ⚠️ 部分模块（net/http、net/ws、net/rpc）需要补充实现
- ⚠️ 工程化相关内容（mod、build）需要补充
- ⚠️ 缺少数据库操作相关示例

### 推荐行动顺序

1. **立即修复**：3 个测试失败（1-2 小时）
2. **短期补充**：net/http 示例 + mod 文档（6-9 小时）
3. **中期完善**：database + 并发模式 + 常量模块（12-15 小时）
4. **长期规划**：crypto + fuzzing + ws/rpc（10-15 小时）

---

## 八、参考资源

### 官方文档

- Go 官方文档：<https://go.dev/doc/>
- Go Blog：<https://go.dev/blog/>
- Effective Go：<https://go.dev/doc/effective_go>
- Go by Example：<https://gobyexample.com/>

### 学习推荐

- 先学完 `builtin/` 再学 `concurrency/`
- `net/http/` 学完后可了解 `net/tcp/`、`net/udp/` 底层原理
- 并发编程建议配合 GMP 调度器文档深入学习

---

*报告状态：Oracle 审查已完成 ✅*
*最后更新：2026-02-26*

---

## 九、待办事项清单 (Todo List)

> **注**：基于代码库深度分析生成，包含具体文件路径、代码位置和实现细节

### P0 - 紧急修复（阻断性问题）

#### 1. 修复测试失败 ❌

| 文件                                      | 行号 | 问题                                             | 修复方案                                 |
| ----------------------------------------- | ---- | ------------------------------------------------ | ---------------------------------------- |
| `builtin/reflects/values/value_test.go`   | 85   | `assert.True(t, tv.CanSet())` 对值类型返回 false | 改为 `assert.False(t, tv.CanSet())`      |
| `runtime/callerstate/callerstate_test.go` | 70   | 硬编码行号断言，期望 68 但实际为 69              | 改为 `assert.Equal(t, 69, cs[0].LineNo)` |

**说明**：

- `TestStrings_HasSuffix` 测试实际通过（strings 包测试全部 ok），无需修复
- `TestReflect_GetValue` 失败原因：`reflect.ValueOf(obj)` 返回值类型，`CanSet()` 恒为 false
- `TestCallerState_ListStackInfo` 失败原因：`ListStackInfo(10)` 调用在第 69 行，行号应为 69

---

### P1 - 高优先级

#### 2. 补充 net/http 模块实现

**当前状态**：仅有 `http_test.go` 测试文件，使用标准库

| 任务              | 文件路径                 | 参考模式            | 核心内容                             |
| ----------------- | ------------------------ | ------------------- | ------------------------------------ |
| 创建协议定义      | `net/http/proto.go`      | `net/tcp/proto.go`  | Request/Response 结构体、Header 定义 |
| 创建日志模块      | `net/http/log.go`        | `net/tcp/log.go`    | sLog、cLog 日志对象                  |
| 创建服务端        | `net/http/server.go`     | `net/tcp/server.go` | Server 结构体、路由、Handler         |
| 创建客户端        | `net/http/client.go`     | `net/tcp/client.go` | Client 结构体、Request、超时控制     |
| 添加中间件示例    | `net/http/middleware.go` | -                   | Logging、Recovery、CORS 中间件       |
| 添加 RESTful 示例 | `net/http/restful.go`    | -                   | RESTful API 设计模式                 |

#### 3. 完善 mod 模块文档

**当前状态**：`mod/README.md` 仅有标题 `# Go Module`

| 章节       | 内容要点                                                               |
| ---------- | ---------------------------------------------------------------------- |
| 模块初始化 | `go mod init`、模块路径命名规范                                        |
| 依赖管理   | `go get`、`go mod tidy`、`go mod verify`、`go mod download`            |
| 版本控制   | 语义化版本、伪版本 (v0.0.0-yyyymmddhhmmss-abcdefabcdef)、indirect 依赖 |
| 私有模块   | `GOPRIVATE`、`GONOSUMDB`、`GOPROXY` 配置                               |
| 工作区     | `go.work`、多模块管理、`go.work.use` 指令                              |
| 常用命令   | `go mod graph`、`go mod why`、`go mod edit`                            |

---

### P2 - 中优先级

#### 4. 创建 database 目录

**当前状态**：不存在，需新建

| 子模块            | 文件             | 内容                                  |
| ----------------- | ---------------- | ------------------------------------- |
| `database/sql/`   | `connect.go`     | 连接池配置、DSN 格式                  |
|                   | `crud.go`        | CRUD 操作示例（Query、Exec、Prepare） |
|                   | `transaction.go` | 事务处理（Begin、Commit、Rollback）   |
|                   | `sql_test.go`    | 完整测试覆盖                          |
| `database/redis/` | `connect.go`     | Redis 连接、连接池配置                |
|                   | `operations.go`  | 基本操作（Set、Get、Del、Expire）     |
|                   | `redis_test.go`  | 测试文件                              |

#### 5. 创建 concurrency/patterns 目录

**当前状态**：不存在，需新建

**已有模式**（无需重复）：

- Worker Pool：`concurrency/sync/pools/worker_pool/worker_pool.go` ✅
- Generator：`concurrency/goroutine/chans/chan.go` ✅
- 阻塞队列：`concurrency/sync/blockque/blockque.go` ✅

| 需补充模式     | 文件路径              | 核心实现                                     |
| -------------- | --------------------- | -------------------------------------------- |
| Fan-out/Fan-in | `patterns/fanout/`    | 多 worker 从同一 channel 读取，结果合并      |
| Pipeline       | `patterns/pipeline/`  | 多阶段处理链（stage-by-stage）               |
| Pub/Sub        | `patterns/pubsub/`    | 发布订阅模式                                 |
| ErrGroup       | `patterns/errgroup/`  | 并发错误处理（`golang.org/x/sync/errgroup`） |
| Rate Limiting  | `patterns/ratelimit/` | 令牌桶、漏桶限流                             |
| Context 组合   | `patterns/context/`   | WithCancel/WithTimeout/WithValue 组合示例    |

#### 6. 创建 builtin/constants 模块

**当前状态**：不存在，需新建

| 文件                | 内容                                  |
| ------------------- | ------------------------------------- |
| `constants.go`      | 常量定义基础                          |
| `iota.go`           | iota 用法示例（枚举、位运算、表达式） |
| `untyped.go`        | 无类型常量说明                        |
| `constants_test.go` | 测试文件                              |

#### 7. 创建 builtin/control 模块

**当前状态**：不存在，需新建

| 文件              | 内容                                            |
| ----------------- | ----------------------------------------------- |
| `for.go`          | for 循环详解（标准、while 风格、死循环、range） |
| `if.go`           | if/else 语句（初始化语句、作用域）              |
| `switch.go`       | switch 语句（表达式、类型、无表达式）           |
| `select.go`       | select 多路复用（超时、非阻塞、默认）           |
| `control_test.go` | 测试文件                                        |

---

### P3 - 低优先级

#### 8. 创建 crypto 目录

**当前状态**：不存在，需新建

| 子模块               | 文件      | 内容                           |
| -------------------- | --------- | ------------------------------ |
| `crypto/hash/`       | `hash.go` | MD5、SHA256、SHA512 哈希       |
| `crypto/symmetric/`  | `aes.go`  | AES 对称加密/解密              |
| `crypto/asymmetric/` | `rsa.go`  | RSA 非对称加密/解密、签名/验证 |

#### 9. 补充 testing/benchmark 模块

| 文件                | 内容                                     |
| ------------------- | ---------------------------------------- |
| `benchmark_test.go` | 基准测试编写、`b.N` 使用                 |
| `compare_test.go`   | 性能比较（`ReportAllocs`、`ResetTimer`） |
| `memory_test.go`    | 内存分配分析（`AllocsPerRun`）           |

#### 10. 补充 builtin/generics 深入

**当前状态**：`builtin/types/generic/` 已有基础内容

| 补充内容     | 文件             |
| ------------ | ---------------- |
| 类型约束详解 | `constraints.go` |
| 类型推断原理 | `inference.go`   |
| 泛型方法限制 | `methods.go`     |

---

### P4 - 待规划

#### 11. 实现 net/ws WebSocket 模块

**当前状态**：`net/ws/ws.go` 为空文件（仅 `package ws`）

| 文件        | 内容               |
| ----------- | ------------------ |
| `conn.go`   | WebSocket 连接封装 |
| `server.go` | WebSocket 服务端   |
| `client.go` | WebSocket 客户端   |
| `frame.go`  | 帧解析（RFC 6455） |

#### 12. 实现 net/rpc RPC 模块

**当前状态**：`net/rpc/rpc.go` 为空文件（仅 `package rpc`）

| 文件        | 内容       |
| ----------- | ---------- |
| `server.go` | RPC 服务端 |
| `client.go` | RPC 客户端 |
| `proto.go`  | 协议定义   |

#### 13. 补充 io/serialize 实现文件

**当前状态**：仅有测试文件 `json_test.go`、`xml_test.go`，缺少实现文件

| 文件      | 内容                                                    |
| --------- | ------------------------------------------------------- |
| `json.go` | JSON 编解码示例（Marshal、Unmarshal、Encoder、Decoder） |
| `xml.go`  | XML 编解码示例                                          |

#### 14. 创建 build 目录

| 文件               | 内容                               |
| ------------------ | ---------------------------------- |
| `build.go`         | go build 命令详解                  |
| `ldflags.go`       | ldflags 编译参数（版本注入、瘦身） |
| `cross_compile.go` | 交叉编译（GOOS、GOARCH）           |

#### 15. 创建 codegen 目录

| 文件          | 内容              |
| ------------- | ----------------- |
| `generate.go` | go generate 命令  |
| `stringer.go` | stringer 工具使用 |

---

## 十、模块文件组织规范

### 代码风格参考

基于现有代码库分析，新模块应遵循以下规范：

| 规范项   | 规则            | 示例                                     |
| -------- | --------------- | ---------------------------------------- |
| 包名     | 小写简写        | `slices`, `strings`, `types`             |
| 函数名   | PascalCase      | `PtrAdd`, `Range`, `SwapElement`         |
| 泛型参数 | 大写 T          | `[T any]`, `[T ~int]`, `[S ~[]T, T any]` |
| 测试函数 | Test 前缀       | `TestSlice_Define`, `TestGeneric_Add`    |
| 测试包名 | `_test` 后缀    | `package slice_test`                     |
| 文件名   | 小写下划线      | `utils.go`, `slice_test.go`              |
| 导入路径 | 项目路径/模块名 | `study/basic/builtin/slices`             |

### 测试风格

```go
package slice_test

import (
    "testing"
    slices2 "study/basic/builtin/slices"
    "github.com/stretchr/testify/assert"
)

func TestSlice_Define(t *testing.T) {
    // 使用 testify/assert 做断言
    assert.Nil(t, s)
    assert.Equal(t, 5, len(s))
}
```

### net 模块结构参考

```
net/{module}/
├── server.go      # 服务端实现
├── client.go      # 客户端实现
├── proto.go       # 协议数据结构
├── log.go         # 日志对象
├── conn.go        # 连接封装（可选）
└── {module}_test.go  # 测试文件
```

---

> **更新时间**：2026-02-26
> **完成进度**：0 / 60+ 项
> **数据来源**：代码库深度探索 + 并行分析
