## Install and setup go environment

### Install GO

#### macOS

```sh
$ brew install go
```



### Setup GOROOT and GOPATH

Open profile file (like `.bash_profile` on macOS)

```sh
export GOROOT=dir				# each of golang binary installed path
export GOPATH=dir1:dir2:dir3	# all go project path
```

### Install glide

#### macOS

```sh
$ brew install golide
```



## Create go project

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

