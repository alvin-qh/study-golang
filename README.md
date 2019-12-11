# Install and setup go environment

## Install GO

### macOS

```sh
$ brew install go
```

## Setup GOROOT and GOPATH

Open profile file (like `.bash_profile` on macOS)

```sh
export GOROOT=dir				# each of golang binary installed path
export GOPATH=dir1:dir2:dir3	# all go project path
```

## Install glide

### macOS

```sh
$ brew install glide
```

## Install xcode-select to enable debug

```bash
$ xcode-select --install
```

## Create GO project

In each GOPATH, create necessary dirs

```sh
$ cd ~/<gopath dir>
$ mkdir src
$ mkdir pkg
$ mkdir bin
```

Create project dir in GOPATH

```sh
$ mkdir myproject
```

Run `glide create`, install depend

```sh
$ glide create		# create project
$ glide get "<go package url>"	# download dependency packages
$ glide up			# upgrade dependency packages
$ glide install		# lock versions of dependency packages
```

## Use GO Module

### Setup

Get version of golang

```bash
$ go version
```

When GO version upper than `1.11` and less than `1.13`

```bash
$ export GO111MODULE=auto
```

or

```bash
$ export GO111MODULE=on
```

### Create GO module

```bash
$ mkdir demo-work
$ cat demo-work
$ go mod init demo-work
```

### Get thirdpart library

```bash
$ go get -u github.com/stretchr/testify/
```

### Set GO proxy

```bash
export GOPROXY=https://goproxy.io
```



