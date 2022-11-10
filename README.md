# 安装和设置 Go 开发环境

- [安装和设置 Go 开发环境](#安装和设置-go-开发环境)
  - [1. 安装 GO 开发环境](#1-安装-go-开发环境)
    - [1.1. macOS](#11-macos)
    - [1.2. Linux](#12-linux)
  - [~~2. 设置 `GOROOT` 和 `GOPATH` 环境变量~~](#2-设置-goroot-和-gopath-环境变量)
  - [~~3. `glide` 包管理器~~](#3-glide-包管理器)
    - [~~3.1. macOS~~](#31-macos)
    - [~~3.2. 通过 `glide` 创建工程~~](#32-通过-glide-创建工程)
  - [4. 使用 `go mod` 包管理](#4-使用-go-mod-包管理)
    - [4.1. 启用包管理](#41-启用包管理)
    - [4.2. 创建 go 工程](#42-创建-go-工程)
    - [4.3. 安装第三方依赖包](#43-安装第三方依赖包)
    - [4.4. 设置 `go mod` 代理](#44-设置-go-mod-代理)

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

## ~~2. 设置 `GOROOT` 和 `GOPATH` 环境变量~~

```bash
export GOROOT=<dir> # 设置为 go 语言安装路径
export GOPATH=<dir1>:<dir2>:<dir3> # 设置为各个 go 工程路径
```

## ~~3. `glide` 包管理器~~

### ~~3.1. macOS~~

安装 `glide` 软件包

```bash
brew install glide
```

安装 `xcode-select`

```bash
xcode-select --install
```

### ~~3.2. 通过 `glide` 创建工程~~

在 `GOPATH` 环境变量指定的路径下, 创建如下目录

```bash
mkdir src  # 源代码路径
mkdir pkg  # 第三方包存放路径
mkdir bin  # 编译结果存放路径
```

在 `GOPATH` 路径下创建工程目录

```bash
mkdir my_project
```

通过 `glide` 创建项目

```bash
glide create                  # 创建工程
glide get "<go package url>"  # 安装第三方依赖包
glide up                      # 更新依赖包
glide install                 # 锁定包版本
```

## 4. 使用 `go mod` 包管理

### 4.1. 启用包管理

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

### 4.2. 创建 go 工程

```bash
mkdir demo-work

cd demo-work
go mod init demo-work
```

### 4.3. 安装第三方依赖包

```bash
go get -u github.com/stretchr/testify/
```

### 4.4. 设置 `go mod` 代理

```bash
export GOPROXY=https://goproxy.io
```
