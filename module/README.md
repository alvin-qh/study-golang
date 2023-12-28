# Go 模块

- [Go 模块](#go-模块)
  - [1. 项目内模块依赖](#1-项目内模块依赖)
  - [2. 引用模块依赖](#2-引用模块依赖)
    - [2.1. 建立依赖模块的代码仓库](#21-建立依赖模块的代码仓库)
    - [2.2. 为依赖模块设置版本号](#22-为依赖模块设置版本号)
    - [2.3. 使用模块](#23-使用模块)
      - [2.3.1. 使用公共仓库](#231-使用公共仓库)
      - [2.3.2. 使用私有仓库](#232-使用私有仓库)
        - [标记指定地址为私有仓库地址](#标记指定地址为私有仓库地址)
        - [设置私有仓库的用户名密码](#设置私有仓库的用户名密码)
        - [仓库地址不合规处理](#仓库地址不合规处理)
      - [2.3.3. 使用本地依赖](#233-使用本地依赖)
    - [2.4. 依赖不一致](#24-依赖不一致)

Go 语言通过"模块"来组织代码, 一个 Go 项目即一个"模块", 根据"模块"的组织形式, 又可以分为:

1. 项目内模块, 即程序和其依赖在同一个模块下, 依赖作为"子模块"存在;
2. 引用模块, 即程序的依赖是一个独立的模块, 并托管在 github 仓库中;
3. 本地引用模块, 即程序的依赖是一个独立模块, 且位于本地磁盘中;

下面详细介绍如何使用上述三种依赖方式

假设主程序模块名称为 `study-golang/module` (定义位于主程序模块目录下的 `go.mod` 文件中), 对于不同类型子模块引用, 有如下说明

## 1. 项目内模块依赖

假设子模块名为 `sub1`

这是最基本和简单的依赖方式, 即依赖在同一个模块的子目录 (`sub1`) 下, 具有自己的模块名 (`package sub1`), 则只需要通过当前模块的名称加上子模块名即可完成引用, 即: `import study-golang/module/sub1` 即可进行引用

## 2. 引用模块依赖

假设子模块名为 `gitee.com/go-common-libs/sub1`, 位于 <https://gitee.com/go-common-libs/sub1> 地址

由于 Go 语言默认的依赖仓库是 github.com, 所以用本例提供的 gitee.com 存放依赖模块, 则需要进行额外的设置, 详细方法可以参考 [在 Git 上托管模块](../module/module1/README.md) 一节的内容

### 2.1. 建立依赖模块的代码仓库

在 gitee.com 上建立路径为 `/go-common-libs/sub1` 的代码仓库并 clone 到本地, 在其中编写所需代码

在编写代码的过程中, 可以正常进行各类 git 操作

### 2.2. 为依赖模块设置版本号

代码编写完成后, 通过 git tag 方式为代码设置版本号

```bash
git tag -a "v1.0.0" -m "tag comment"
git push --tags
```

注意, 有效分支必须为主分支 (`master` 或 `main`) 分支

### 2.3. 使用模块

根据代码仓库的可见性 (公共或私有) 不同, 作为依赖的处理方式也略有不同

#### 2.3.1. 使用公共仓库

如果要使用的模块位于 git 的一个公共仓库 (Public Repository), 则直接使用 `go get` 命令下载安装即可, 即:

```bash
go get -u gitee.com/go-common-libs/sub1@v1.0.0
```

此时会在 `go.mod` 文件中加入如下记录, 表示依赖安装成功

```plaintext
require (
    gitee.com/go-common-libs/sub1 v1.0.0 // indirect
)
```

之后即可在代码中使用 `import gitee.com/go-common-libs/sub1/xxx` 来引用依赖和其中的子模块

#### 2.3.2. 使用私有仓库

使用私有仓库的难点在于私有仓库中的代码需要提供用户凭证 (即 git 的用户名和密码), 需要如下额外的步骤

##### 标记指定地址为私有仓库地址

通过环境变量 `GOPRIVATE` 可以标记指定的 git 地址为私有仓库, 此时 Go 在获取该仓库代码作为依赖时, 会忽略一些安全检查, 例如:

```bash
export GOPRIVATE="gitee.com"
```

上述命令表示将整个 gitee.com 下的仓库均设置为私有仓库, Go 将不再对 gitee.com 下的代码仓库进行检查

也可以只设置部分地址, 例如:

```bash
export GOPRIVATE="gitee.com/go-libs/sub1,gitee.com/go-libs/sub2"
```

当然, 如果 Go 版本大于 1.13, 则可以用 `go env -w` 命令来设置 Go 环境变量, 更为方便

```bash
go env -w GOPRIVATE="gitee.com/go-libs/sub1,gitee.com/go-libs/sub2"
```

##### 设置私有仓库的用户名密码

要成功拉取依赖代码, 还需要设置代码仓库的访问权限, 即用户名密码

编辑 `~/.netrc` 文件, 为指定的 git 地址设置用户名和密码, 内容可以如下:

```plaintext
machine gitee.com
login quhao317@163.com
password **************
```

如果不希望在文件中记录用户名和密码, 也可以按如下方式设置环境变量, 这样会在每次获取模块 (执行 `go get`) 时提示密码输入

```bash
export GIT_TERMINAL_PROMPT=1
```

完成上述配置后, 即可像访问公共仓库一样使用私有仓库

除了 `GOPRIVATE` 变量外, 还有其它变量可以辅助 Go 使用私有仓库中的模块:

- `GONOPROXY`: 指定不使用默认代理的仓库地址, 例如: `go env -w GONOPROXY="gitee.com/go-common-libs"`
- `GONOSUMDB`: 指定不进行校验 (`go sum`) 的仓库地址, 例如: `go env -w GONOSUMDB="gitee.com/go-common-libs"`
- `GOINSECURE`: 指定不使用 `https` 的仓库地址, 例如: `go env -w GOINSECURE="gitee.com/go-common-libs"`
  - 也可以通过 `go get --insecure gitee.com/go-common-lib/<repository-name>` 在每次拉取模块时使用, 但比较麻烦

##### 仓库地址不合规处理

Go 默认使用 github.com 仓库进行模块管理, 所以不允许模块所在仓库的地址中包含数字和端口号, 此时可以通过 git 的一些设置进行处理

假设模块的仓库地址为 <https://git.my-private.com:9090/go-libs/sub1>, 则可以对该地址进行 URL 替换设置

```bash
git config --global url."https://git.my-private.com:9090/".insteadof "https://git.my-private-git.com/"
go env -w GOPRIVATE="*.my-private-git.com"
```

完成设置后, 所有访问到 git.my-private-git.com 的访问都会替换为 git.my-private.com:9090, 从而绕过了 Go 的 URL 检查规则, 此时如下命令可以正常执行

```bash
go get -u git.my-private-git.com/go-libs/sub1
```

#### 2.3.3. 使用本地依赖

在开发时, 如果依赖库还不成熟, 需要频繁改动, 则不断推送到 git 仓库比较麻烦, 或者依赖模块的代码本身并未在 git 仓库托管, 此时可以使用本地磁盘目录作为依赖仓库

在主程序模块中添加本地依赖模块, 需要在引入模块 (执行 `go get` 命令) 前, 在主程序模块的 `go.mod` 文件中增加一行说明

```plaintext
replace study-golang/module/sub1 => ../module2
```

此说明表示将要引入一个名为 `study-golang/module/sub1` 的模块, 其位置位于 `../module2` 目录下

再次说明的基础上, 即可为主程序安装依赖模块

```bash
go get -v -u study-golang/module/sub1
```

此时会返回 `go: added study-golang/module/sub1 v0.0.0-00010101000000-000000000000` 表示添加成功, 但因为目标代码不在 git 仓库中, 没有表示版本的 tag, 所以随机给了一个版本号. 如果觉得这个随即版本号比较奇怪, 可以在主程序模块的 `go.mod` 文件中将其改为任意版本号, 不影响正常执行

### 2.4. 依赖不一致

如果在获取依赖是报告 "go module xxxx found but does not contain package", 一般是因为远端依赖的 go 代码发生了改动, 但并未升级版本号, 导致本地缓存和远端代码不一致, 此时清空本地缓存即可

```bash
go clean -cache
go clean -modcache
```
