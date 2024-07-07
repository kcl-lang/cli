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

const (
	VersionTypeLatest = Version_0_9_1

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
