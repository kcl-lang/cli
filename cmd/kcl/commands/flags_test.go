// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
	"kcl-lang.io/cli/pkg/options"
)

// TestAppendLangFlags_CommaSafeValues verifies that every repeatable flag
// registered by appendLangFlags/appendRunnerFlags keeps the value of a
// single occurrence intact even when the value contains commas, brackets
// or quoted substrings. Previously the flags used pflag.StringSliceVarP,
// which splits values via encoding/csv and breaks legitimate inputs such
// as `-O person.ids=[1,2]` and `-O 'person.name="Alice,Bob"'`.
// Regression test for https://github.com/kcl-lang/cli/issues/232.
func TestAppendLangFlags_CommaSafeValues(t *testing.T) {
	cases := []struct {
		name string
		args []string
		// collect returns the slice stored on RunOptions for the flag
		// under test. Using accessors keeps the test independent of the
		// field ordering inside RunOptions.
		collect func(*options.RunOptions) []string
		want    []string
	}{
		{
			name:    "overrides list value with commas",
			args:    []string{"-O", "person.ids=[1,2,3]"},
			collect: func(o *options.RunOptions) []string { return o.Overrides },
			want:    []string{"person.ids=[1,2,3]"},
		},
		{
			name:    "overrides string value with quoted comma",
			args:    []string{`-O`, `person.name="Alice,Bob"`},
			collect: func(o *options.RunOptions) []string { return o.Overrides },
			want:    []string{`person.name="Alice,Bob"`},
		},
		{
			name: "overrides multiple occurrences",
			args: []string{
				"-O", "person.age=10",
				"-O", "person.name=\"Bob\"",
				"-O", "person.ids=[1,2]",
			},
			collect: func(o *options.RunOptions) []string { return o.Overrides },
			want:    []string{`person.age=10`, `person.name="Bob"`, `person.ids=[1,2]`},
		},
		{
			name:    "path selector keeps commas",
			args:    []string{"-S", "a.b.c"},
			collect: func(o *options.RunOptions) []string { return o.PathSelectors },
			want:    []string{"a.b.c"},
		},
		{
			name:    "settings keeps commas",
			args:    []string{"-Y", "settings.yaml,with,comma.yaml"},
			collect: func(o *options.RunOptions) []string { return o.Settings },
			want:    []string{"settings.yaml,with,comma.yaml"},
		},
		{
			name:    "external packages keeps commas",
			args:    []string{"-E", "my_pkg=./vendor/my,pkg"},
			collect: func(o *options.RunOptions) []string { return o.ExternalPackages },
			want:    []string{"my_pkg=./vendor/my,pkg"},
		},
		{
			name: "argument list value with commas",
			args: []string{"-D", "ids=[1,2,3]"},
			collect: func(o *options.RunOptions) []string {
				return o.Arguments
			},
			want: []string{"ids=[1,2,3]"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := options.NewRunOptions()
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			appendLangFlags(o, flags)
			if err := flags.Parse(tc.args); err != nil {
				t.Fatalf("Parse(%v) returned error: %v", tc.args, err)
			}
			got := tc.collect(o)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Parse(%v) = %#v, want %#v", tc.args, got, tc.want)
			}
		})
	}
}
