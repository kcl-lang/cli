# FreeBSD

The CLI loads the KCL runtime from `libkcl.so` and calls `kcl_service_new`.

## Startup

On FreeBSD, commands that evaluate KCL (`run`, `lint`, `fmt`, `test`, `vet`, `server`, and `doc generate`) look for `libkcl.so` in this order:

1. `$KCL_LIB_HOME/libkcl.so`
2. `$XDG_CACHE_HOME/kcl/libkcl.so`, or `~/.cache/kcl/libkcl.so` when `XDG_CACHE_HOME` is unset
3. `/usr/local/lib/libkcl.so` (the FreeBSD ports and packages prefix)
4. `/usr/lib/libkcl.so`

A usable library exports `kcl_service_new`. The CLI then sets `KCL_LIB_HOME` to that library's directory so `kcl-lang.io/lib` opens the same file. A library found under `/usr/local/lib` or `/usr/lib` is linked into the cache directory, which is the path the loader already uses.

`libkclvm_cli_cdylib.so` is an older artifact. The loader does not open it, and a file that only exports `kclvm_service_new` is rejected.

Before the native client starts, the CLI calls `install.InstallKcl` from `kcl-lang.io/lib` on that directory. A FreeBSD `installLib` can unpack the runtime into the cache when nothing usable is on disk yet. This repository does not implement `installLib`. If the runtime is still missing, or the file does not export `kcl_service_new`, the command prints an error that names `libkcl` and `freebsd/<arch>` and exits. It does not panic with `Undefined symbol "kclvm_service_new"`.

`kcl version` and `kcl --help` do not load the runtime.

## Building

GitHub-hosted runners do not provide a FreeBSD host. Cross-compile from Linux or another supported Go host:

```shell
make build-freebsd
```

The equivalent command for one architecture is:

```shell
GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build \
  -gcflags='github.com/ebitengine/purego/internal/fakecgo=-std' \
  -o bin/kcl-freebsd-amd64 ./cmd/kcl
```

`github.com/ebitengine/purego` needs that `-gcflags` value when `CGO_ENABLED=0` or when cross-compiling. It makes fakecgo define `environ` and `__progname`, which FreeBSD libc expects. `amd64` and `arm64` are both built by `make build-freebsd`.

## Remaining release gap

`kcl-lang.io/lib` v0.13.2 builds `go/install/install_lib_unix.go` for Linux and Darwin only. A FreeBSD compile therefore stops in the dependency:

```text
kcl-lang.io/lib/go/install/install.go: undefined: installLib
```

CI runs the cross-compile and treats that specific error as a known gap (`ALLOW_FREEBSD_LIB_GAP=1`). Any other compiler error still fails the job. `.goreleaser.yml` keeps release archives on Darwin, Linux, and Windows until FreeBSD `libkcl` artifacts and `installLib` exist.

Runtime artifacts for FreeBSD are tracked in [kcl-lang/kcl#1029](https://github.com/kcl-lang/kcl/issues/1029). The undefined-symbol panic is [kcl-lang/cli#147](https://github.com/kcl-lang/cli/issues/147).
