# 安装和设置 Go 开发环境

- [安装和设置 Go 开发环境](#安装和设置-go-开发环境)
  - [1. 安装 GO 开发环境](#1-安装-go-开发环境)
    - [1.1. macOS](#11-macos)
    - [1.2. Linux](#12-linux)
  - [3. 使用 go mod 包管理](#3-使用-go-mod-包管理)
    - [3.1. 启用包管理](#31-启用包管理)
    - [3.2. 创建 go 工程](#32-创建-go-工程)
    - [3.3. 安装第三方依赖包](#33-安装第三方依赖包)
    - [3.4. 设置 `go mod` 代理](#34-设置-go-mod-代理)
  - [4. 使用 go env 环境变量管理](#4-使用-go-env-环境变量管理)

![Logo](./assets/logo.jpg)

## 1. 安装 GO 开发环境

可以从 [Download and install - The Go Programming Language (golang.org)](https://golang.org/doc/install) 下载二进制安装包直接进行安装, 除此之外, 各平台也提供了通过各自管理工具进行安装的方法

### 1.1. macOS

通过 brew 管理工具进行安装

```bash
brew install go
```

### 1.2. Linux

添加软件源

```bash
sudo add-apt-repository ppa:longsleep/golang-backports
```

将软件源替换为国内代理

```bash
sudo vi /etc/apt/sources.list.d/longsleep-ubuntu-golang-backports-focal.list
```

```bash
# deb http://ppa.launchpad.net/longsleep/golang-backports/ubuntu jammy main
deb https://launchpad.proxy.ustclug.org/longsleep/golang-backports/ubuntu jammy main
```

安装

```bash
sudo apt update -y
sudo apt install golang-go
```

## 3. 使用 go mod 包管理

### 3.1. 启用包管理

查看 go 版本号

```bash
go version
```

对于版本号高于 `1.11` 的, 启用 `go mod` 模块

```bash
export GO111MODULE=auto
```

或者

```bash
export GO111MODULE=on
```

### 3.2. 创建 go 工程

```bash
mkdir demo-work

cd demo-work
go mod init gitee.com/go-libs/demo-work
```

### 3.3. 安装第三方依赖包

```bash
go get -u github.com/stretchr/testify/
```

### 3.4. 设置 `go mod` 代理

```bash
export GOPROXY="https://goproxy.cn,direct"
```

## 4. 使用 go env 环境变量管理

对于上述需要设置环境变量的环节, 如果 Go 版本大于 `1.13` 则可以直接使用 `go env -w` 命令

```bash
go env -w GO111MODULE="on"
go env -w GOPROXY="https://goproxy.cn,direct"
```

`go env -w` 设置的环境变量存储在 `~/.config/go/env` 文件中 (Linux 系统下), 并且可以跨平台使用

通过 `go env` 命令可以列举目前所有 Go 环境变量的值
