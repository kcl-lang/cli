package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kpm/pkg/client"
)

// TestParseSourceFromArgs_HostlessOci tests that a bare repo path passed to
// --oci (no scheme, no host) produces an Oci source with RegFromEnv=true,
// preserving the host-less form so KPM_REG resolves the registry at runtime.
func TestParseSourceFromArgs_HostlessOci(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	// Simulate: kcl mod add --oci myorg/kcl-templates/utils --tag 0.2.0
	oci = "myorg/kcl-templates/utils"
	tag = "0.2.0"
	defer func() { oci = ""; tag = "" }()

	src, err := ParseSourceFromArgs(kpmcli, []string{})
	assert.NoError(t, err)
	assert.NotNil(t, src.Oci, "expected an OCI source")
	assert.Equal(t, "myorg/kcl-templates/utils", src.Oci.Repo)
	assert.Equal(t, "0.2.0", src.Oci.Tag)
	assert.Equal(t, "", src.Oci.Reg, "registry must be empty (resolved at runtime from KPM_REG)")
	assert.True(t, src.Oci.RegFromEnv, "RegFromEnv must be set for host-less deps")
}

// TestParseSourceFromArgs_HostnameWithoutScheme tests the behaviour when --oci
// receives a value that looks like it contains a registry host but has no URL
// scheme (e.g. "ghcr.io/kcl-lang/helloworld").  Go's url.Parse cannot detect a
// host without an explicit scheme, so the entire string is treated as a bare
// repository path — the same code path as a path-only input — and RegFromEnv is
// set to true.  "ghcr.io" is NOT extracted as the registry host.
//
// This test documents that edge-case so the behaviour is explicit.  If the
// intent ever changes (i.e. the code is taught to recognise "host/path" without
// a scheme), this test must be updated accordingly.
func TestParseSourceFromArgs_HostnameWithoutScheme(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	// Simulate: kcl mod add --oci ghcr.io/kcl-lang/helloworld --tag 0.1.0
	oci = "ghcr.io/kcl-lang/helloworld"
	tag = "0.1.0"
	defer func() { oci = ""; tag = "" }()

	src, err := ParseSourceFromArgs(kpmcli, []string{})
	assert.NoError(t, err)
	assert.NotNil(t, src.Oci, "expected an OCI source")
	// url.Parse sees no scheme → the whole string becomes the path; "ghcr.io"
	// is NOT separated out as the registry host.
	assert.Equal(t, "ghcr.io/kcl-lang/helloworld", src.Oci.Repo)
	assert.Equal(t, "", src.Oci.Reg, "registry must be empty: host-like prefix without scheme is not parsed as a host")
	assert.True(t, src.Oci.RegFromEnv, "host-like prefix without scheme triggers RegFromEnv like a bare path")
}

// TestParseSourceFromArgs_FullOciUrl tests that a full OCI URL (oci://host/repo)
// produces a source with Reg populated and RegFromEnv=false (legacy behaviour).
func TestParseSourceFromArgs_FullOciUrl(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	// Simulate: kcl mod add --oci oci://ghcr.io/kcl-lang/helloworld --tag 0.1.0
	oci = "oci://ghcr.io/kcl-lang/helloworld"
	tag = "0.1.0"
	defer func() { oci = ""; tag = "" }()

	src, err := ParseSourceFromArgs(kpmcli, []string{})
	assert.NoError(t, err)
	assert.NotNil(t, src.Oci, "expected an OCI source")
	assert.Equal(t, "ghcr.io", src.Oci.Reg)
	assert.Equal(t, "kcl-lang/helloworld", src.Oci.Repo)
	assert.Equal(t, "0.1.0", src.Oci.Tag)
	assert.False(t, src.Oci.RegFromEnv, "RegFromEnv must be false for full-URL deps")
}

// TestParseSourceFromArgs_HostlessOciWithModSpec tests that a ModSpec positional
// argument (e.g. "mypkg:1.0") can be combined with a host-less --oci flag.
// The resulting Source must carry both the OCI coordinates (with RegFromEnv=true)
// and the ModSpec name/version so KPM can locate the right sub-module inside the
// registry image.
func TestParseSourceFromArgs_HostlessOciWithModSpec(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	// Simulate: kcl mod add mypkg:1.0 --oci myorg/kcl-templates/utils --tag 0.2.0
	oci = "myorg/kcl-templates/utils"
	tag = "0.2.0"
	defer func() { oci = ""; tag = "" }()

	src, err := ParseSourceFromArgs(kpmcli, []string{"mypkg:1.0"})
	assert.NoError(t, err)

	assert.NotNil(t, src.Oci, "expected an OCI source")
	assert.Equal(t, "myorg/kcl-templates/utils", src.Oci.Repo)
	assert.Equal(t, "0.2.0", src.Oci.Tag)
	assert.True(t, src.Oci.RegFromEnv, "RegFromEnv must be set for host-less deps")

	assert.NotNil(t, src.ModSpec, "expected a ModSpec from the positional argument")
	assert.Equal(t, "mypkg", src.ModSpec.Name)
	assert.Equal(t, "1.0", src.ModSpec.Version)
}
