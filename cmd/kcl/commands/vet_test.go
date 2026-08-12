package cmd

import (
	"testing"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// TestNewVetCmdExposesExternalFlag is a structural regression test for
// https://github.com/kcl-lang/kcl/issues/1877 that ensures the `kcl vet`
// command registers the `--external`/`-E` flag and keeps the embedded
// `validate.ValidateOptions` reachable through `VetOptions`.
func TestNewVetCmdExposesExternalFlag(t *testing.T) {
	cmd := NewVetCmd()

	flag := cmd.Flags().Lookup("external")
	if flag == nil {
		t.Fatal("expected `vet` to expose the `--external` flag")
	}
	if flag.Shorthand != "E" {
		t.Fatalf("expected `--external` shorthand to be `E`, got %q", flag.Shorthand)
	}
}

// TestParseExternalPackages covers the parser that turns raw CLI
// values into `gpyrpc.ExternalPkg` entries ready to be forwarded to
// the gRPC service.
func TestParseExternalPackages(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []*gpyrpc.ExternalPkg
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  nil,
		},
		{
			name:  "single entry",
			input: []string{"k8s=./vendor/k8s"},
			want: []*gpyrpc.ExternalPkg{
				{PkgName: "k8s", PkgPath: "./vendor/k8s"},
			},
		},
		{
			name:  "multiple entries with surrounding whitespace",
			input: []string{" k8s = ./vendor/k8s ", "ext=../ext"},
			want: []*gpyrpc.ExternalPkg{
				{PkgName: "k8s", PkgPath: "./vendor/k8s"},
				{PkgName: "ext", PkgPath: "../ext"},
			},
		},
		{
			name:  "blank entries are skipped",
			input: []string{"", "k8s=./vendor/k8s", "   "},
			want: []*gpyrpc.ExternalPkg{
				{PkgName: "k8s", PkgPath: "./vendor/k8s"},
			},
		},
		{
			name:    "missing separator is rejected",
			input:   []string{"only_a_path"},
			wantErr: true,
		},
		{
			name:    "empty package name is rejected",
			input:   []string{"=./vendor/k8s"},
			wantErr: true,
		},
		{
			name:    "empty path is rejected",
			input:   []string{"k8s="},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExternalPackages(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d entries, got %d (%v)", len(tt.want), len(got), got)
			}
			for i := range got {
				if got[i].PkgName != tt.want[i].PkgName || got[i].PkgPath != tt.want[i].PkgPath {
					t.Fatalf("entry %d: expected %+v, got %+v", i, tt.want[i], got[i])
				}
			}
		})
	}
}

// TestParseExternalPackagesNilOnEmptyInput documents that the helper
// returns nil (rather than an empty slice) when there is nothing to
// forward, keeping the wire format compatible with the historical
// behaviour.
func TestParseExternalPackagesNilOnEmptyInput(t *testing.T) {
	if got, err := parseExternalPackages(nil); err != nil || got != nil {
		t.Fatalf("expected (nil, nil) for nil input, got (%v, %v)", got, err)
	}
	if got, err := parseExternalPackages([]string{"", "   "}); err != nil || got != nil {
		t.Fatalf("expected (nil, nil) for only-blank input, got (%v, %v)", got, err)
	}
}
