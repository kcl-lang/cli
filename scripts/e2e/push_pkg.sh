
#!/usr/bin/env bash

export KPM_REG="localhost:5001"
export KPM_REPO="test"
export OCI_REG_PLAIN_HTTP=on

# Prepare the package on the registry
current_dir=$(pwd)
echo $current_dir

$current_dir/bin/kcl registry login -u test -p 1234 localhost:5001

cd ./scripts/e2e/pkg_in_reg/oci/ghcr.io/kcl-lang/k8s/1.27/k8s/1.27
$current_dir/bin/kcl mod push

cd "$current_dir"

cd ./scripts/e2e/pkg_in_reg/oci/ghcr.io/kcl-lang/k8s/1.28/k8s/1.28
$current_dir/bin/kcl mod push

cd "$current_dir"

cd ./scripts/e2e/pkg_in_reg/oci/ghcr.io/kcl-lang/k8s/1.31.2/k8s/1.31.2
$current_dir/bin/kcl mod push

cd "$current_dir"

# Push the package helloworld/0.1.1 to the registry
cd ./scripts/e2e/pkg_in_reg/oci/ghcr.io/kcl-lang/helloworld/0.1.1/helloworld/0.1.1
$current_dir/bin/kcl mod push

cd "$current_dir"

# Push the package helloworld/0.1.4 to the registry
cd ./scripts/e2e/pkg_in_reg/oci/ghcr.io/kcl-lang/helloworld/0.1.4/helloworld/0.1.4
$current_dir/bin/kcl mod push

cd "$current_dir"

# Push the package 'kcl1' depends on 'k8s' to the registry
cd ./scripts/e2e/pkg_in_reg/kcl1
$current_dir/bin/kcl mod push

cd "$current_dir"

# Push the package 'kcl2' depends on 'k8s' to the registry
cd ./scripts/e2e/pkg_in_reg/kcl2
$current_dir/bin/kcl mod push
