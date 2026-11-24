package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/downloader"
	"kcl-lang.io/kpm/pkg/opt"
	"kcl-lang.io/kpm/pkg/utils"
)

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

func ParseSourceFromArgs(cli *client.KpmClient, args []string) (*downloader.Source, error) {
	source := downloader.Source{}
	modSpec := downloader.ModSpec{
		Alias: rename,
	}

	// Parse the source from the args
	// Parse the input like: kcl mod pull k8s:1.28 or kcl mod pull oci://ghcr.io/kcl-lang/helloworld --tag 0.1.0
	if len(args) != 0 {
		// Iterate through the args to find the source
		for _, arg := range args {
			// if arg is a path and exists, it is a local source
			if utils.DirExists(arg) {
				source.Local = &downloader.Local{
					Path: arg,
				}
				continue
			}

			// A reference that names the registry or the namespace cannot be
			// expressed by a package spec, so it is parsed as an OCI source.
			if len(git) == 0 && len(oci) == 0 && len(path) == 0 {
				if ociSource, ok := parseQualifiedOCIRef(arg); ok {
					if source.Oci != nil {
						return nil, fmt.Errorf("only one source is allowed")
					}
					source.Oci = ociSource
					continue
				}
			}

			// if arg is not path, it is Modspec
			if err := modSpec.FromString(arg); err == nil {
				// check if modspec already exists
				if source.ModSpec == nil {
					source.ModSpec = &modSpec
				} else {
					return nil, fmt.Errorf("only one modspec is allowed")
				}
				continue
			} else {
				modSpec = downloader.ModSpec{
					Alias: rename,
				}
			}

			// if arg is a url, set the source url
			if source.IsNilSource() {
				err := source.FromString(arg)
				if err != nil {
					return nil, err
				}
				continue
			}

			if !source.IsNilSource() {
				return nil, fmt.Errorf("only one source is allowed")
			}
		}
		// For the source parsed from the args, set the tag, commit, branch
		if source.Git != nil {
			source.Git.Tag = tag
			source.Git.Commit = commit
			source.Git.Branch = branch
		} else if source.Oci != nil && len(tag) != 0 {
			// The tag flag overrides the tag of the reference, but it must
			// not clear it: the reference may carry a tag itself, either
			// in the "name:tag" form or in the query of an "oci://" URL.
			source.Oci.Tag = tag
		}
	}

	// Parse the source from the flags
	// Parse the input like: kcl mod pull --oci oci://ghcr.io/kcl-lang/helloworld --tag 0.1.0
	if source.IsNilSource() || source.SpecOnly() {
		if len(git) != 0 {
			source.Git = &downloader.Git{
				Url:    git,
				Tag:    tag,
				Commit: commit,
				Branch: branch,
			}
		} else if len(oci) != 0 {
			ociUrl, err := url.Parse(oci)
			if err != nil {
				return nil, err
			}

			// If the value has no scheme at all, treat it as a bare repository path
			// (e.g. "myorg/kcl-templates/utils").  The registry host will be resolved
			// at runtime from KPM_REG / DefaultOciRegistry.
			if ociUrl.Scheme == "" {
				source.Oci = &downloader.Oci{
					Repo:       oci,
					Tag:        tag,
					RegFromEnv: true,
				}
			} else {
				ociUrl.Scheme = constants.OciScheme
				query := ociUrl.Query()
				query.Add(constants.Tag, tag)
				ociUrl.RawQuery = query.Encode()
				err = source.FromString(ociUrl.String())
				if err != nil {
					return nil, err
				}
			}
		} else if len(path) != 0 {
			source.Local = &downloader.Local{
				Path: path,
			}
		}
	} else if len(git) != 0 || len(oci) != 0 || len(path) != 0 {
		return nil, fmt.Errorf("only one source is allowed")
	}

	source.ModSpec = &modSpec

	return &source, nil
}

// parseQualifiedOCIRef parses a package reference that names the registry
// or the namespace explicitly:
//
//	<namespace>/<name>[:<tag>]
//	<registry>/<namespace>/<name>[:<tag>]
//
// It returns the OCI source and true when the argument is such a
// reference, and false for anything else — a bare package name, a URL, a
// local path or a git reference.
//
// Such a reference cannot be expressed by a package spec: the spec's name
// is joined with the default repo when it is resolved, which turns
// "ghcr.io/kcl-lang/helloworld" into "kcl-lang/ghcr.io/kcl-lang/helloworld"
// and "kcl-lang/helloworld" into "kcl-lang/kcl-lang/helloworld" — both
// repositories that do not exist.
//
// The namespace replaces the default repo instead of being prefixed to it,
// the same way the host-less --oci form treats a bare repository path. The
// registry stays unresolved when the reference does not name one, so it is
// read from KPM_REG at runtime.
func parseQualifiedOCIRef(arg string) (*downloader.Oci, bool) {
	// A URL is parsed as a URL elsewhere, and a digest reference carries an
	// "@" that is not part of a tag.
	if strings.Contains(arg, "://") || strings.Contains(arg, "@") {
		return nil, false
	}

	// The tag is the part after the last colon that follows the last slash,
	// so that a port in the registry host ("localhost:5000/foo") is not
	// mistaken for a tag.
	name, tag := arg, ""
	if idx := strings.LastIndex(arg, ":"); idx != -1 && !strings.Contains(arg[idx+1:], "/") {
		name, tag = arg[:idx], arg[idx+1:]
	}

	segments := strings.Split(name, "/")
	if len(segments) < 2 {
		// A bare package name lives in the default repo.
		return nil, false
	}

	reg := ""
	if isRegistryHost(segments[0]) {
		reg = segments[0]
		segments = segments[1:]
	}
	if len(segments) == 0 {
		return nil, false
	}

	return &downloader.Oci{
		Reg:  reg,
		Repo: strings.Join(segments, "/"),
		Tag:  tag,
		// The registry is resolved at runtime from KPM_REG when the
		// reference does not name one.
		RegFromEnv: reg == "",
	}, true
}

// isRegistryHost reports whether the first segment of an OCI reference
// names a registry rather than a namespace.
//
// It follows the convention the OCI distribution spec and Docker use: the
// first segment is a registry host when it contains a dot or a colon, or
// when it is exactly "localhost". Anything else is a namespace.
func isRegistryHost(segment string) bool {
	if segment == "localhost" {
		return true
	}
	return strings.ContainsAny(segment, ".:")
}

func ParseUrlFromArgs(cli *client.KpmClient, args []string) (*url.URL, error) {
	var sourceUrl url.URL

	if len(args) == 0 {
		if len(git) != 0 {
			gitUrl, err := url.Parse(git)
			if err != nil {
				return nil, err
			}

			gitUrl.Scheme = constants.GitScheme
			query := gitUrl.Query()
			if tag != "" {
				query.Add(constants.Tag, tag)
			}
			if commit != "" {
				query.Add(constants.GitCommit, commit)
			}
			if branch != "" {
				query.Add(constants.GitBranch, branch)
			}
			gitUrl.RawQuery = query.Encode()
			sourceUrl = *gitUrl
		} else if len(oci) != 0 {
			ociUrl, err := url.Parse(oci)
			if err != nil {
				return nil, err
			}

			ociUrl.Scheme = constants.OciScheme
			query := ociUrl.Query()
			query.Add(constants.Tag, tag)
			ociUrl.RawQuery = query.Encode()
			sourceUrl = *ociUrl
		}
	} else {
		url, err := url.Parse(args[0])
		if err != nil {
			return nil, err
		}
		query := url.Query()
		url.Opaque = ""
		regOpts, err := opt.NewRegistryOptionsFrom(args[0], cli.GetSettings())
		if err != nil {
			return nil, err
		}

		if regOpts.Git != nil {
			if url.Scheme != constants.GitScheme && url.Scheme != constants.SshScheme {
				url.Scheme = constants.GitScheme
			}
			if tag != "" {
				query.Add(constants.Tag, tag)
			}
			if commit != "" {
				query.Add(constants.GitCommit, commit)
			}
			if branch != "" {
				query.Add(constants.GitBranch, branch)
			}
		} else if regOpts.Oci != nil {
			url.Scheme = constants.OciScheme
			url.Host = regOpts.Oci.Reg
			url.Path = regOpts.Oci.Repo
			if regOpts.Oci.Tag != "" {
				query.Add(constants.Tag, regOpts.Oci.Tag)
			}
			if tag != "" {
				query.Add(constants.Tag, tag)
			}
		} else if regOpts.Registry != nil {
			url.Scheme = constants.DefaultOciScheme
			url.Host = regOpts.Registry.Reg
			url.Path = regOpts.Registry.Repo
			if regOpts.Registry.Tag != "" {
				query.Add(constants.Tag, regOpts.Registry.Tag)
			}
			if tag != "" {
				query.Add(constants.Tag, tag)
			}
		}

		url.RawQuery = query.Encode()
		sourceUrl = *url
	}
	return &sourceUrl, nil
}
