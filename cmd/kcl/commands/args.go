package cmd

import (
	"net/url"

	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/opt"
)

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
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
