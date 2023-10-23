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

`kcl` 是一个命令行界面，包括语言核心功能、IDE功能、包管理工具、社区集成和其他工具等。

## 安装

### 使用脚本安装

### 使用 `go install` 安装

您可以使用 `go install` 命令安装 `kcl`。

```shell
go install kcl-lang.io/cli@latest
```

### 从 Github Release 页面手动安装

您也可以从 Github Release 中获取 `kcl` ，并将 `kcl` 的二进制文件路径设置到环境变量 PATH 中。

```shell
# KCL_INSTALLATION_PATH 是 `kcl` 二进制文件的所在目录.
export PATH=$KCL_INSTALLATION_PATH:$PATH  
```

请使用以下命令以确保您成功安装了 `kcl`。

```shell
kcl --help
```

### 从源代码构建

## 快速开始

```shell
kcl run ./examples/kubernetes.k
```

## 常见问题 (FAQ)

##### Q: 我在使用 `go install` 安装 `kcl` 后，出现了 `command not found` 的错误。

A: `go install` 默认会将二进制文件安装到 `$GOPATH/bin` 目录下，您需要将 `$GOPATH/bin` 添加到环境变量 `PATH` 中。
