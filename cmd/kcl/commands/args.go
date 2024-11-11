package cmd

import (
	"fmt"
	"net/url"

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
				modSpec = downloader.ModSpec{}
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
		} else if source.Oci != nil {
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

			ociUrl.Scheme = constants.OciScheme
			query := ociUrl.Query()
			query.Add(constants.Tag, tag)
			ociUrl.RawQuery = query.Encode()
			err = source.FromString(ociUrl.String())
			if err != nil {
				return nil, err
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
