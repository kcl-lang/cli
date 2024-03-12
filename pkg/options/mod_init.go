package options

import "kcl-lang.io/kpm/pkg/errors"

// Input options of 'kpm init'.
type ModInitOptions struct {
	Name     string
	InitPath string
}

func (opts *ModInitOptions) Validate() error {
	if len(opts.Name) == 0 {
		return errors.InvalidInitOptions
	} else if len(opts.InitPath) == 0 {
		return errors.InternalBug
	}
	return nil
}