// Copyright The KCL Authors. All rights reserved.

package options

import (
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHasLintPattern(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{args: nil, want: false},
		{args: []string{"main.k"}, want: false},
		{args: []string{"pkg/...x", "main.k"}, want: false},
		{args: []string{"./..."}, want: true},
		{args: []string{"pkg/..."}, want: true},
		{args: []string{"..."}, want: true},
		{args: []string{"main.k", "./..."}, want: true},
	}
	for _, tt := range tests {
		assert.Equal(t, HasLintPattern(tt.args), tt.want, "args: %v", tt.args)
	}
}

func TestLintPackages(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    [][]string
		wantErr bool
	}{
		{
			name: "pattern selects every package recursively",
			args: []string{"./testdata/lint/..."},
			want: [][]string{
				{"testdata/lint/main.k"},
				{"testdata/lint/nested/sub/warn.k"},
				{"testdata/lint/pkg_dup/dup.k"},
			},
		},
		{
			name: "plain directory selects its root package only",
			args: []string{"./testdata/lint"},
			want: [][]string{{"testdata/lint/main.k"}},
		},
		{
			name: "plain file joins its directory package",
			args: []string{"./testdata/lint/pkg_dup/dup.k", "./testdata/lint/..."},
			want: [][]string{
				{"testdata/lint/main.k"},
				{"testdata/lint/nested/sub/warn.k"},
				{"testdata/lint/pkg_dup/dup.k"},
			},
		},
		{
			name:    "no KCL files is an error",
			args:    []string{filepath.Join(t.TempDir(), "...")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lintPackages(tt.args)
			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}
			assert.NilError(t, err)
			assert.DeepEqual(t, got, tt.want)
		})
	}
}

// TestLintAllPackagesLintsEveryPackageIndependently lints a tree where the
// schema `S` is declared by two distinct packages: linting every package as
// an independent program does not collide the duplicate declarations, while
// linting both files as one program does.
func TestLintAllPackagesLintsEveryPackageIndependently(t *testing.T) {
	// Both files in a single program redeclare the schema `S`.
	merged := NewRunOptions()
	merged.Entries = []string{"./testdata/lint/main.k", "./testdata/lint/pkg_dup/dup.k"}
	merged.CompileOnly = true
	assert.Assert(t, merged.Run() != nil)

	// Lint every package independently instead.
	o := NewRunOptions()
	o.NoStyle = true
	err := o.LintAllPackages([]string{"./testdata/lint/main.k", "./testdata/lint/pkg_dup"})
	assert.NilError(t, err)
}

// TestLintAllPackagesReportsFindings ensures lint findings of a package fail
// the lint and are attributed to the package directory.
func TestLintAllPackagesReportsFindings(t *testing.T) {
	o := NewRunOptions()
	o.NoStyle = true
	err := o.LintAllPackages([]string{"./testdata/lint/nested/..."})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, filepath.Join("testdata", "lint", "nested", "sub"))
	assert.ErrorContains(t, err, "Module 'math' imported but unused")
}
