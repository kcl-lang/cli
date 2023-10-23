// Copyright The KCL Authors. All rights reserved.

package main

import (
	"os"

	"kcl-lang.io/cli/cmd"

	_ "kcl-lang.io/kcl-go"
	_ "kcl-lang.io/kcl-openapi/pkg/kube_resource/generator"
	_ "kcl-lang.io/kpm/pkg/oci"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
