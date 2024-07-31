#!/bin/sh

# Stop on error.

# Note: Before you start tests, please install kcl
# kcl Installation: https://kcl-lang.io/docs/user_docs/getting-started/install

set -e

pwd=$(
    cd $(dirname $0)
    pwd
)

for path in "configuration" "validation" "abstraction" "definition" "konfig" "mutation" "data-integration" "automation" "package-management" "kubernetes" "codelab" "server" "settings" "source"; do
    echo "\033[1mTesting $path ...\033[0m"
    if (cd $pwd/$path && make test); then
        echo "\033[32mTest SUCCESSED - $path\033[0m\n"
    else
        echo "\033[31mTest FAILED - $path\033[0m\n"
        exit 1
    fi
done
