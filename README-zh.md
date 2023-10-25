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

`kcl` 是一个命令行界面，包括语言核心功能、IDE 功能、包管理工具、社区集成和其他工具等。

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
git clone https://github.com/kcl-lang/cli
cd cli && go build ./cmd/kcl/main.go -o kcl
```

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
