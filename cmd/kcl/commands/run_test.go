package cmd

import (
	"strings"
	"testing"
)

func TestNewRunCmd(t *testing.T) {
	cmd := NewRunCmd()

	if cmd.Use != "run" {
		t.Errorf("unexpected command use: %s", cmd.Use)
	}

	if cmd.Short != "Run KCL codes." {
		t.Errorf("unexpected command short description: %s", cmd.Short)
	}

	if cmd.Long != runDesc {
		t.Errorf("unexpected command long description: %s", cmd.Long)
	}

	if cmd.Example != runExample {
		t.Errorf("unexpected command example: %s", cmd.Example)
	}

	if cmd.SilenceUsage != true {
		t.Errorf("unexpected SilenceUsage value: %v", cmd.SilenceUsage)
	}

	runE := cmd.RunE
	if runE == nil {
		t.Fatal("RunE function is nil")
	}

	args := []string{"../../../examples/kubernetes.k"}
	err := runE(cmd, args)
	if err != nil {
		t.Errorf("RunE function returned an error: %v", err)
	}

	args = []string{"error.k"}
	err = runE(cmd, args)
	if !strings.Contains(err.Error(), "Cannot find the kcl file") {
		t.Errorf("RunE function returned an error: %v", err)
	}

	args = []string{"error.k"}
	err = runE(cmd, args)
	if !strings.Contains(err.Error(), "Cannot find the kcl file") {
		t.Errorf("RunE function returned an error: %v", err)
	}
}
