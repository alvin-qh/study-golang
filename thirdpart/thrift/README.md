# Apache Thrift

## 1. 安装 CLI 工具

### 1.1. Windows 系统

Thrift 为 Windows 系统提供了编译过后的 `.exe` 文件, 可以直接下载运行, 下载连接为 <https://thrift.apache.org/download>

也可以通过 Scoop 等包管理器直接安装

```bash
scoop install thrift
```

### 1.2. macOS 安装

macOS 可以通过 brew 包管理器直接安装

```bash
brew install thrift
```

### 1.3. Linux 安装

Linux 需要下载源代码并进行编译, 下面以 Debian 系统为例

#### 1.3.1. 下载源代码

从 <https://thrift.apache.org/download> 链接下载源代码, 并解压到任意目录

#### 1.3.2. 安装工具链和依赖库

安装工具链

```bash
sudo apt install build-essential g++ autotools-dev
```

安装依赖库

```bash
sudo apt install python-dev-is-python3 \
     libicu-dev \
     libevent-dev \
     zlib1g-dev \
     libbz2-dev \
     libboost-all-dev \
     libssl-dev
```

安装其它工具

```bash
sudo apt install libtool \
     flex \
     bison
```

#### 1.3.3. 执行编译

如果需要多开发环境支持, 则需要为各开发环境安装支持, 包括:

- Java: 需安装 Ant (`sudo apt install ant`) 及 Gradle
- NodeJS: 需要安装 node 及 npm 包管理器, 也可以使用 nvm
- Python: 需要 Python 3.8 以上版本和 SetupTool (`python -m pip install setuptools`)
- Golang: 需安装 Go 工具链
- Rust: 需安装 rustup, 并通过 rustup 安装 cargo 及 rustc 工具链
- C/C++: 需安装 gcc 和 g++

以上开发环境需安装并配置正确的环境变量

由于部分编译命令需要在 root 用户权限下执行, 所以需要做如下操作:

```bash
# 如果通过 nvm 安装 node, 则需要如下操作
sudo ln -s "$(which node)" /usr/bin/node
sudo ln -s "$(which npm)" /usr/bin/npm

# 如果需要 Rust 支持, 则需要如下操作
sudo su -
cd ~
cp -r /home/<user>/.cargo .  # 将 .cargo 目录复制到 root 下
rustup default stable
```

进入 thrift 源码目录, 执行如下操作:

```bash
./bootstrap.sh
autoupdate

./configure

make -j 10
make install
```
