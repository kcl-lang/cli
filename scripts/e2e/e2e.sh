#!/usr/bin/env bash

# start registry at 'localhost:5001'
# include account 'test' and password '1234'
./scripts/e2e/reg.sh

# set the kpm default registry and repository
export KPM_REG="localhost:5001"
export KPM_REPO="test"
export OCI_REG_PLAIN_HTTP=on
# kpm 0.12.9 caches remote primary packages referenced by `kcl run`
# (oci://... / --git / --oci) by default. Disable that here so the
# existing e2e goldens — which were captured against the always-fresh
# pre-cache behaviour — keep matching. The new caching feature is
# exercised in kpm's own test suite; these cli tests are about
# kcl-run output, not kpm cache state.
export KPM_RUN_NO_CACHE=1

set -o errexit
set -o nounset
set -o pipefail

# Install ginkgo
GO111MODULE=on go install github.com/onsi/ginkgo/v2/ginkgo@v2.0.0

# Build kpm binary
make build

# Prepare e2e test env
# pull the package 'k8s' from 'ghcr.io/kcl-lang/k8s'
./scripts/e2e/pull_pkg.sh

# push the package 'k8s' to 'localhost:5001/test'
./scripts/e2e/push_pkg.sh

# Run e2e
set +e
ginkgo  ./test/e2e/ 
TESTING_RESULT=$?


exit $TESTING_RESULT
