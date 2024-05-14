#!/bin/bash

concurrent_test() {
    local repo_url=$1
    local concurrency=$2

    local -a statuses
    for i in $(seq 1 "$concurrency"); do
        statuses[i]=0
    done

    run_test() {
        local id=$1
        echo "Starting test for $repo_url (ID: $id)"
        if time kcl run "$repo_url"; then
            echo "Completed test successfully for $repo_url (ID: $id)"
            statuses[id]=0
        else
            echo "Test failed for $repo_url (ID: $id)"
            statuses[id]=$?
        fi
    }

    for i in $(seq 1 "$concurrency"); do
        run_test "$i" &
    done

    wait

    local has_errors=0
    for status in "${statuses[@]}"; do
        if [ "$status" -ne 0 ]; then
            has_errors=1
            break
        fi
    done

    return $has_errors
}

TEST_REPOS=(
    "oci://ghcr.io/kcl-lang/podinfo"
    "https://github.com/kcl-lang/flask-demo-kcl-manifests"
    "./examples/server"
)
CONCURRENCY_LEVEL=4

for repo in "${TEST_REPOS[@]}"; do
    if ! concurrent_test "$repo" "$CONCURRENCY_LEVEL"; then
        echo "Error during concurrent test for $repo"
        exit 1
    fi
done

echo "All concurrent tests completed successfully."
