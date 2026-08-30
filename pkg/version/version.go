// Copyright The KCL Authors. All rights reserved.

package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

// version will be set by build flags.
var version string

// GetVersionString() returns the KCL CLI version string in the form
// `{version}-{goos}-{goarch}`.
//
// The version is resolved in the following order:
//  1. The `version` package variable, populated by `-ldflags "-X .../pkg/version.version=..."`
//     at build time (used by goreleaser and the musl release build).
//  2. The module version embedded in the Go build info (works for binaries installed
//     via `go install kcl-lang.io/cli@vX.Y.Z` from a tagged module).
//  3. The hardcoded `VersionTypeLatest` constant as a last resort.
func GetVersionString() string {
	if len(version) != 0 {
		return version
	}
	if v := buildInfoVersion(); v != "" {
		return v
	}
	return VersionTypeLatest.String()
}

// buildInfoVersion extracts a usable version string from the Go runtime build info.
// It strips the leading "v" (e.g. "v0.12.8" -> "0.12.8") and returns an empty
// string when the build info is unavailable or reports a non-release version
// such as "(devel)" or a pseudo-version.
func buildInfoVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	v := info.Main.Version
	if v == "" || v == "(devel)" {
		return ""
	}
	return strings.TrimPrefix(v, "v")
}

// VersionType is the version type of kpm.
type VersionType string

// String() will transform VersionType to string.
func (kvt VersionType) String() string {
	return getVersion(string(kvt))
}

func getVersion(version string) string {
	return fmt.Sprintf("%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)
}

const (
	VersionTypeLatest = Version_0_12_10

	Version_0_12_10 VersionType = "0.12.10"
	Version_0_12_9 VersionType = "0.12.9"
	Version_0_12_8 VersionType = "0.12.8"
	Version_0_12_7 VersionType = "0.12.7"
	Version_0_12_6 VersionType = "0.12.6"
	Version_0_12_5 VersionType = "0.12.5"
	Version_0_12_4 VersionType = "0.12.4"
	Version_0_12_3 VersionType = "0.12.3"
	Version_0_12_2 VersionType = "0.12.2"
	Version_0_12_1 VersionType = "0.12.1"
	Version_0_12_0 VersionType = "0.12.0"
	Version_0_11_2 VersionType = "0.11.2"
	Version_0_11_1 VersionType = "0.11.1"
	Version_0_11_0 VersionType = "0.11.0"
	Version_0_10_9 VersionType = "0.10.9"
	Version_0_10_8 VersionType = "0.10.8"
	Version_0_10_7 VersionType = "0.10.7"
	Version_0_10_6 VersionType = "0.10.6"
	Version_0_10_5 VersionType = "0.10.5"
	Version_0_10_4 VersionType = "0.10.4"
	Version_0_10_3 VersionType = "0.10.3"
	Version_0_10_2 VersionType = "0.10.2"
	Version_0_10_1 VersionType = "0.10.1"
	Version_0_10_0 VersionType = "0.10.0"

	Version_0_9_8 VersionType = "0.9.8"
	Version_0_9_7 VersionType = "0.9.7"
	Version_0_9_6 VersionType = "0.9.6"
	Version_0_9_5 VersionType = "0.9.5"
	Version_0_9_4 VersionType = "0.9.4"
	Version_0_9_3 VersionType = "0.9.3"
	Version_0_9_2 VersionType = "0.9.2"
	Version_0_9_1 VersionType = "0.9.1"
	Version_0_9_0 VersionType = "0.9.0"

	Version_0_8_9 VersionType = "0.8.9"
	Version_0_8_8 VersionType = "0.8.8"
	Version_0_8_7 VersionType = "0.8.7"
	Version_0_8_6 VersionType = "0.8.6"
	Version_0_8_5 VersionType = "0.8.5"
	Version_0_8_4 VersionType = "0.8.4"
	Version_0_8_3 VersionType = "0.8.3"
	Version_0_8_2 VersionType = "0.8.2"
	Version_0_8_1 VersionType = "0.8.1"
	Version_0_8_0 VersionType = "0.8.0"

	Version_0_7_5 VersionType = "0.7.5"
	Version_0_7_4 VersionType = "0.7.4"
	Version_0_7_3 VersionType = "0.7.3"
	Version_0_7_2 VersionType = "0.7.2"
	Version_0_7_1 VersionType = "0.7.1"
	Version_0_7_0 VersionType = "0.7.0"
	Version_0_6_0 VersionType = "0.6.0"
)
