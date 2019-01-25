# Api Backend

![Build Status](https://img.shields.io/travis/rust-lang/rust.svg)


## Contents

- [SETUP](#markdown-header-setup)
- [GIT PRE HOOK SETUP](#markdown-header-git-pre-hook-setup)

## SETUP

* Download and install go, dep ( Go dependency management tool )
* Install mockery (https://github.com/vektra/mockery)

* Set GOPATH locally for Go workspace and add in your bash

```sh
export GOPATH=/Users/{{name}}/your_folder_path
export PATH=$PATH:$GOPATH/bin
```

* Form dir structure according to your project repo

```sh
cd $GOPATH
mkdir -p src/bitbucket.org/libertywireless
```

* Go to created dir and clone the project

```sh
cd src/bitbucket.org/libertywireless
git clone git@bitbucket.org:libertywireless/go-restapi-boilerplate.git
```

* Go inside the project and create log file

```sh
cd go-restapi-boilerplate
mkdir -p log/rest_api.log
```

* Install project dependencies and build

```go
go get
go build
```

* Execute the project executable file created using above command.

```go
 ./go-restapi-boilerplate
```

* Go to localhost:3030/health and project setup is done.

## GIT PRE HOOK SETUP

* Install pre-commit package manager

```sh
brew install pre-commit
```

* Run install command to install pre-commit into your git hooks

```sh
pre-commit install
```

* Install go lint package using command

```go
 go get -u golang.org/x/lint/golint
```

* pre-commit will now run on every commit.
