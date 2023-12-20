// Copyright The KCL Authors. All rights reserved.

package version

import (
	"fmt"
	"runtime"
)

// version will be set by build flags.
var version string

// GetVersionString() will return the latest version of kpm.
func GetVersionString() string {
	if len(version) == 0 {
		// If version is not set by build flags, return the version constant.
		return VersionTypeLatest.String()
	}
	return version
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

// All the kpm versions.
const (
	VersionTypeLatest = Version_0_7_3_alpha_1

	Version_0_7_3_alpha_1 VersionType = "0.7.3-alpha.1"
	Version_0_7_2         VersionType = "0.7.2"
	Version_0_7_1         VersionType = "0.7.1"
	Version_0_7_0         VersionType = "0.7.0"
	Version_0_7_0_beta_2  VersionType = "0.7.0-beta.2"
	Version_0_7_0_beta_1  VersionType = "0.7.0-beta.1"
	Version_0_7_0_alpha_2 VersionType = "0.7.0-alpha.2"
	Version_0_7_0_alpha_1 VersionType = "0.7.0-alpha.1"
	Version_0_6_0         VersionType = "0.6.0"
	Version_0_6_0_alpha_1 VersionType = "0.6.0-alpha.1"
	Version_0_6_0_alpha_2 VersionType = "0.6.0-alpha.2"
	Version_0_6_0_alpha_3 VersionType = "0.6.0-alpha.3"
	Version_0_6_0_alpha_4 VersionType = "0.6.0-alpha.4"
)
