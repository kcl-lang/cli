<h1 align="center">KCL CLI</h1>

<p align="center">
<a href="./README.md">English</a> | <a href="./README-zh.md">简体中文</a>
</p>
<p align="center">
<a href="#introduction">Introduction</a> | <a href="#installation">Installation</a> | <a href="#quick-start">Quick start</a> 
</p>

<p align="center">
<img src="https://coveralls.io/repos/github/kcl-lang/cli/badge.svg">
<img src="https://img.shields.io/badge/license-Apache--2.0-green">
<img src="https://img.shields.io/badge/PRs-welcome-brightgreen">
<img src="https://img.shields.io/github/downloads/kcl-lang/cli/total?label=Github%20downloads&logo=github">
</p>

## Introduction
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkcl-lang%2Fcli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkcl-lang%2Fcli?ref=badge_shield)

KCL `cli` is a command-line interface that includes KCL language core features, IDE features, package management tools, plugins, community integration, and other tools. It now integrates the following tools:

+ [KCL Core and API](https://github.com/kcl-lang/kcl)
+ [KCL Go SDK](https://github.com/kcl-lang/kcl-go)
+ [KCL Open API Tool](https://github.com/kcl-lang/kcl-openapi)
+ [KCL Package Manage Tool](https://github.com/kcl-lang/kpm)
+ [KCL Plugins](https://github.com/kcl-lang/kcl-plugin)
+ [KCL Playground](https://github.com/kcl-lang/kcl-playground)

## Installation

### Scripts

#### MacOS

```shell
curl -fsSL https://kcl-lang.io/script/install-cli.sh | /bin/bash
```

#### Linux

```shell
wget -q https://kcl-lang.io/script/install-cli.sh -O - | /bin/bash
```

#### Windows

```shell
powershell -Command "iwr -useb https://kcl-lang.io/script/install-cli.ps1 | iex"
```

### Homebrew (MacOS & Linux)

```shell
brew install kcl-lang/tap/kcl
```

### Scoop (Windows)

```shell
scoop bucket add kcl-lang https://github.com/kcl-lang/scoop-bucket.git
scoop install kcl-lang/kcl
```

### Go install

You can download `kcl` via `go install`.

```shell
go install kcl-lang.io/cli/cmd/kcl@latest
```

### Download from GITHUB Release Page

You can also get `kcl` from the [github release](https://github.com/kcl-lang/cli/releases) and set the binary path to the environment variable PATH.

```shell
# KCL_CLI_INSTALLATION_PATH is the path of the `KCL CLI` binary.
export PATH=$KCL_CLI_INSTALLATION_PATH:$PATH  
```

### Docker

```shell
docker run -it kcllang/kcl
```

### Docker for arm64

```shell
docker run -it kcllang/kcl-arm64
```

### Build from Source Code

```shell
git clone https://github.com/kcl-lang/cli && cd cli
# On Windows, MacOS and Linux
go build ./cmd/kcl/
# Build on Linux Musl
CGO_ENABLED=1 go build -tags="musl netgo static osusergo" -ldflags="-linkmode external -extldflags '-static'" ./cmd/kcl
```

Use the following command to ensure that you install `kcl` successfully.

```shell
kcl --help
```

## Quick Start

```shell
kcl run ./examples/kubernetes.k
```

## Learn More

- [KCL Website](https://kcl-lang.io)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkcl-lang%2Fcli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkcl-lang%2Fcli?ref=badge_large)
