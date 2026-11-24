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

// TestParseQualifiedOCIRef tests the parsing of a positional reference that
// names the registry or the namespace explicitly.
//
// Such a reference cannot be expressed by a package spec: the spec's name is
// joined with the default repo when it is resolved, which turns
// "ghcr.io/kcl-lang/helloworld" into "kcl-lang/ghcr.io/kcl-lang/helloworld"
// and "kcl-lang/helloworld" into "kcl-lang/kcl-lang/helloworld" — both
// repositories that do not exist.
func TestParseQualifiedOCIRef(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		reg  string
		repo string
		tag  string
		ok   bool
	}{
		{
			name: "registry, namespace and name",
			arg:  "ghcr.io/kcl-lang/helloworld:0.1.0",
			reg:  "ghcr.io",
			repo: "kcl-lang/helloworld",
			tag:  "0.1.0",
			ok:   true,
		},
		{
			name: "registry, namespace and name without a tag",
			arg:  "ghcr.io/kcl-lang/helloworld",
			reg:  "ghcr.io",
			repo: "kcl-lang/helloworld",
			ok:   true,
		},
		{
			name: "namespace and name",
			arg:  "kcl-lang/helloworld:0.1.0",
			reg:  "",
			repo: "kcl-lang/helloworld",
			tag:  "0.1.0",
			ok:   true,
		},
		{
			name: "nested package path",
			arg:  "myorg/kcl-templates/utils:0.2.0",
			reg:  "",
			repo: "myorg/kcl-templates/utils",
			tag:  "0.2.0",
			ok:   true,
		},
		{
			name: "registry with a port",
			arg:  "localhost:5000/myorg/helloworld:0.1.0",
			reg:  "localhost:5000",
			repo: "myorg/helloworld",
			tag:  "0.1.0",
			ok:   true,
		},
		{
			name: "a bare package name lives in the default repo",
			arg:  "helloworld:0.1.0",
			ok:   false,
		},
		{
			name: "a URL is not an OCI reference",
			arg:  "oci://ghcr.io/kcl-lang/helloworld",
			ok:   false,
		},
		{
			name: "a git URL is not an OCI reference",
			arg:  "git@github.com:kcl-lang/konfig",
			ok:   false,
		},
		{
			name: "a digest reference is not an OCI reference",
			arg:  "ghcr.io/kcl-lang/helloworld@sha256:0000000000000000000000000000000000000000000000000000000000000000",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ociSource, ok := parseQualifiedOCIRef(tt.arg)
			assert.Equal(t, tt.ok, ok)
			if !tt.ok {
				assert.Nil(t, ociSource)
				return
			}
			assert.Equal(t, tt.reg, ociSource.Reg)
			assert.Equal(t, tt.repo, ociSource.Repo)
			assert.Equal(t, tt.tag, ociSource.Tag)
			assert.Equal(t, tt.reg == "", ociSource.RegFromEnv,
				"the registry must only be resolved at runtime when the reference does not name one")
		})
	}
}

// TestParseSourceFromArgs_QualifiedRef tests that a positional reference
// naming the registry produces an OCI source with the registry, the
// repository and the tag of the reference — and not a package spec, whose
// name would be joined with the default repo.
func TestParseSourceFromArgs_QualifiedRef(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	// Simulate: kcl mod pull ghcr.io/kcl-lang/helloworld:0.1.0
	src, err := ParseSourceFromArgs(kpmcli, []string{"ghcr.io/kcl-lang/helloworld:0.1.0"})
	assert.NoError(t, err)

	assert.NotNil(t, src.Oci, "expected an OCI source")
	assert.Equal(t, "ghcr.io", src.Oci.Reg)
	assert.Equal(t, "kcl-lang/helloworld", src.Oci.Repo)
	assert.Equal(t, "0.1.0", src.Oci.Tag)
	assert.False(t, src.Oci.RegFromEnv, "the reference names the registry")
}

// TestParseSourceFromArgs_QualifiedRefKeepsItsTag tests that the tag of a
// positional reference survives parsing: the --tag flag overrides it, but
// its absence must not clear it.
func TestParseSourceFromArgs_QualifiedRefKeepsItsTag(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	src, err := ParseSourceFromArgs(kpmcli, []string{"ghcr.io/kcl-lang/helloworld:0.1.0"})
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0", src.Oci.Tag)

	// The --tag flag wins over the tag of the reference.
	tag = "0.1.4"
	defer func() { tag = "" }()
	src, err = ParseSourceFromArgs(kpmcli, []string{"ghcr.io/kcl-lang/helloworld:0.1.0"})
	assert.NoError(t, err)
	assert.Equal(t, "0.1.4", src.Oci.Tag)
}

// TestParseSourceFromArgs_OciUrlKeepsItsQueryTag tests that the tag in the
// query of an "oci://" URL is not cleared when no --tag flag is given.
func TestParseSourceFromArgs_OciUrlKeepsItsQueryTag(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	src, err := ParseSourceFromArgs(kpmcli, []string{"oci://ghcr.io/kcl-lang/helloworld?tag=0.1.0"})
	assert.NoError(t, err)

	assert.NotNil(t, src.Oci, "expected an OCI source")
	assert.Equal(t, "ghcr.io", src.Oci.Reg)
	assert.Equal(t, "kcl-lang/helloworld", src.Oci.Repo)
	assert.Equal(t, "0.1.0", src.Oci.Tag)
}

// TestParseSourceFromArgs_BareNameStaysSpecOnly tests that a bare package
// name is still a package spec: it is written to kcl.mod as
// "name = version" and resolved in the default repo.
func TestParseSourceFromArgs_BareNameStaysSpecOnly(t *testing.T) {
	kpmcli, err := client.NewKpmClient()
	assert.NoError(t, err)

	src, err := ParseSourceFromArgs(kpmcli, []string{"helloworld:0.1.0"})
	assert.NoError(t, err)

	assert.Nil(t, src.Oci, "a bare package name must not become an OCI source")
	assert.NotNil(t, src.ModSpec)
	assert.Equal(t, "helloworld", src.ModSpec.Name)
	assert.Equal(t, "0.1.0", src.ModSpec.Version)
	assert.True(t, src.SpecOnly(), "a bare package name stays spec only")
}
