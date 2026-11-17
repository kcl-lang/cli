<h1 align="center">KCL Command Line Interface (CLI)</h1>

<p align="center">
<a href="./README.md">English</a> | <a href="./README-zh.md">简体中文</a>
</p>
<p align="center">
<a href="#介绍">介绍</a> | <a href="#安装">安装</a> | <a href="#快速开始">快速开始</a>
</p>


<p align="center">
<img src="https://coveralls.io/repos/github/kcl-lang/cli/badge.svg">
<img src="https://img.shields.io/badge/license-Apache--2.0-green">
<img src="https://img.shields.io/badge/PRs-welcome-brightgreen">
</p>

## 介绍

KCL `cli` 是一个命令行界面，包括 KCL 语言核心功能、IDE 功能、包管理工具、插件、社区集成和其他工具等。

+ [KCL 核心和 API](https://github.com/kcl-lang/kcl)
+ [KCL Go SDK](https://github.com/kcl-lang/kcl-go)
+ [KCL Open API 工具](https://github.com/kcl-lang/kcl-openapi)
+ [KCL 包管理工具](https://github.com/kcl-lang/kpm)
+ [KCL 插件](https://github.com/kcl-lang/kcl-plugin)
+ [KCL Playground](https://github.com/kcl-lang/kcl-playground)

## 安装

### 使用脚本安装

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

### 使用 `go install` 安装

您可以使用 `go install` 命令安装 `kcl`。

```shell
go install kcl-lang.io/cli/cmd/kcl@latest
```

### 从 Github Release 页面手动安装

您也可以从 [Github Release](https://github.com/kcl-lang/cli/releases) 中获取 `kcl` ，并将 `kcl` 的二进制文件路径设置到环境变量 PATH 中。

```shell
# KCL_INSTALLATION_PATH 是 `kcl` 二进制文件的所在目录.
export PATH=$KCL_INSTALLATION_PATH:$PATH  
```

### 从源代码构建

```shell
git clone https://github.com/kcl-lang/cli && cd cli
# On Windows, MacOS and Linux
go build ./cmd/kcl/
# Build on Linux Musl
CGO_ENABLED=1 go build -tags="musl netgo static osusergo" -ldflags="-linkmode external -extldflags '-static'" ./cmd/kcl
```

#### FreeBSD

在 FreeBSD 上，CLI 加载 `libkcl.so` 并调用 `kcl_service_new`。查找顺序为 `KCL_LIB_HOME`、用户缓存（`~/.cache/kcl`）、`/usr/local/lib` 和 `/usr/lib`。缺少运行时库时会返回指明 `libkcl` 和平台的错误，而不会 panic。交叉编译仍依赖 [kcl-lang/lib](https://github.com/kcl-lang/lib) 中的 `installLib`，详见 [docs/freebsd.md](docs/freebsd.md) 与 [kcl-lang/kcl#1029](https://github.com/kcl-lang/kcl/issues/1029)。

请使用以下命令以确保您成功安装了 `kcl`。

```shell
kcl --help
```

## 快速开始

```shell
kcl run ./examples/kubernetes.k
```

## 更多资源

- [KCL 网站](https://kcl-lang.io)
