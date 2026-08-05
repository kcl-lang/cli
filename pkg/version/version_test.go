package version

import "testing"

func TestGetVersionString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test get version in string",
			want: VersionTypeLatest.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersionString(); got != tt.want {
				t.Errorf(" GetVersionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildInfoVersion(t *testing.T) {
	// Save and restore the package-level version variable so we can exercise the
	// build info path without leaking state across tests.
	prevVersion := version
	defer func() { version = prevVersion }()

	// Force the package-level fallback so buildInfoVersion is consulted.
	version = ""

	// When running under `go test`, ReadBuildInfo returns "(devel)" for the main
	// module, which buildInfoVersion must treat as "no version available".
	if got := buildInfoVersion(); got != "" {
		t.Errorf("buildInfoVersion() under `go test` = %q, want empty string", got)
	}
}
