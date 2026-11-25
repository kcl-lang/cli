#!/bin/sh
# Strip LLVM LTO bitcode from the vendored musl libkcl.a before cgo links.
#
# golang:1.26-alpine moved past the Go that accepted this archive. Current
# Go rejects objects named by #cgo LDFLAGS when they contain .llvmbc or
# .llvmcmd (cmd/go/internal/work/security.go). That bitcode is unused for
# a normal static link. kcl-lang.io/lib v0.13.3 still ships it.
#
# Only the archive for the current GOARCH is rewritten. A cross objcopy
# cannot parse the other architecture's objects.
set -eu

root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$root"

arch=$(go env GOARCH)
case "$arch" in
  amd64) plat=linux-musl-amd64 ;;
  arm64) plat=linux-musl-arm64 ;;
  *)
    echo "strip-musl-lto: unsupported GOARCH $arch" >&2
    exit 1
    ;;
esac

# go list -m leaves Dir empty until the module is downloaded.
go mod download kcl-lang.io/lib >/dev/null
moddir=$(go list -m -f '{{.Dir}}' kcl-lang.io/lib)
if [ -z "$moddir" ] || [ ! -d "$moddir" ]; then
  echo "strip-musl-lto: could not locate kcl-lang.io/lib" >&2
  exit 1
fi
archive="$moddir/go/lib/$plat/libkcl.a"
if [ ! -f "$archive" ]; then
  echo "strip-musl-lto: missing $archive" >&2
  exit 1
fi

# Module cache files are mode 0444 and their directories are 0555.
dir=$(dirname "$archive")
chmod u+w "$dir" "$archive"
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
if ! objcopy --remove-section=.llvmbc --remove-section=.llvmcmd "$archive" "$tmp"; then
  echo "strip-musl-lto: objcopy failed for $archive" >&2
  exit 1
fi
mv "$tmp" "$archive"
trap - EXIT

if grep -a -F -q -e '.llvmbc' -e '.llvmcmd' "$archive"; then
  echo "strip-musl-lto: LTO sections remain in $archive" >&2
  exit 1
fi

echo "stripped .llvmbc and .llvmcmd from $archive"
