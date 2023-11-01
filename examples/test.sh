#!/bin/sh

# Stop on error.

# Note: Before you start tests, please install kcl and kpm
# kcl Installation: https://kcl-lang.io/docs/user_docs/getting-started/install
# kpm Installation: https://kcl-lang.io/docs/user_docs/guides/package-management/installation

set -e

pwd=$(
    cd $(dirname $0)
    pwd
)

for path in "configuration" "validation" "abstraction" "definition" "mutation" "data-integration" "automation" "package-management" "kubernetes" "codelab"; do
    echo "\033[1mTesting $path ...\033[0m"
    if (cd $pwd/$path && make test); then
        echo "\033[32mTest SUCCESSED - $path\033[0m\n"
    else
        echo "\033[31mTest FAILED - $path\033[0m\n"
        exit 1
    fi
done
