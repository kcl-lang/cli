// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandNeedsLibKCL(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{name: "run", want: true},
		{name: "lint", want: true},
		{name: "fmt", want: true},
		{name: "test", want: true},
		{name: "vet", want: true},
		{name: "server", want: true},
		{name: "generate", want: true},
		{name: "version", want: false},
		{name: "clean", want: false},
		{name: "init", want: false},
		{name: "login", want: false},
		{name: "import", want: false},
	}
	for _, tc := range cases {
		c := &cobra.Command{Use: tc.name}
		if got := commandNeedsLibKCL(c); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
	if commandNeedsLibKCL(nil) {
		t.Fatal("nil command should not require libkcl")
	}
}
