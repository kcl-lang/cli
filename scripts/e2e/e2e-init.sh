#!/bin/bash

# Check if a directory name is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: <test suite name>"
    exit 1
fi

# Specify the directory
dir="./test/e2e/test_suites/$1"

# Create the subdirectory if it does not exist
if [ ! -d "$dir" ]; then
  mkdir -p "$dir"
fi

# Create files in the directory

touch "${dir}/input"
echo "stdout" > "${dir}/stdout"
echo "stderr" > "${dir}/stderr"
mkdir -p "${dir}/test_space"


echo "Test suite created successfully in $dir."