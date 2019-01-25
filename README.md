# Golang RestfulAPI Sample

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
mkdir -p src/github.com/kyawmyintthein
```

* Go to created dir and clone the project

```sh
cd src/github.com/kyawmyintthein
git clone git@github.com/kyawmyintthein/golangRestfulAPISample.git
```

* Go inside the project and create log file

```sh
cd golangRestfulAPISample
mkdir -p log/rest_api.log
```

