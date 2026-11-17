#!/usr/bin/env bash
# Cross-compile the kcl CLI for FreeBSD.
#
# GitHub-hosted runners do not provide a FreeBSD host. This script
# cross-compiles from the current machine. Linking still depends on
# kcl-lang.io/lib providing installLib for freebsd; until that lands the
# build stops with "undefined: installLib". This repository does not
# implement that installer.
#
# purego's FreeBSD fakecgo requires
#   -gcflags=github.com/ebitengine/purego/internal/fakecgo=-std
# when CGO_ENABLED=0 or when cross-compiling, so libc sees environ and
# __progname.
#
# Set ALLOW_FREEBSD_LIB_GAP=1 to exit 0 when that known installLib gap is
# the only failure. Any other compiler error still fails the script.

set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

version="$(cat VERSION)"
out_dir="${1:-bin}"
mkdir -p "$out_dir"

gap=0
for arch in amd64 arm64; do
  dest="$out_dir/kcl-freebsd-${arch}"
  log="$(mktemp)"
  set +e
  GOOS=freebsd GOARCH="$arch" CGO_ENABLED=0 go build \
    -gcflags='github.com/ebitengine/purego/internal/fakecgo=-std' \
    -trimpath \
    -ldflags="-X kcl-lang.io/cli/pkg/version.version=${version}" \
    -o "$dest" \
    ./cmd/kcl >"$log" 2>&1
  code=$?
  set -e
  cat "$log"
  if [[ "$code" -eq 0 ]]; then
    echo "built ${dest}"
    rm -f "$log"
    continue
  fi
  rm -f "$dest"
  if grep -q 'undefined: installLib' "$log"; then
    echo "freebsd/${arch}: blocked on kcl-lang.io/lib installLib. See docs/freebsd.md and https://github.com/kcl-lang/kcl/issues/1029"
    gap=1
    rm -f "$log"
    continue
  fi
  rm -f "$log"
  echo "freebsd/${arch}: cross-compile failed" >&2
  exit "$code"
done

if [[ "$gap" -ne 0 ]]; then
  if [[ "${ALLOW_FREEBSD_LIB_GAP:-}" == "1" ]]; then
    echo "FreeBSD binaries were not produced. The remaining gap is installLib in kcl-lang.io/lib; see docs/freebsd.md."
    exit 0
  fi
  echo "FreeBSD binaries were not produced. See docs/freebsd.md." >&2
  exit 2
fi
