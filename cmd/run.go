package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
	"kcl-lang.io/cli/pkg/options"
)

// NewRunCmd returns the run command.
func NewRunCmd() *cobra.Command {
	o := options.NewRunOptions()
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run KCL codes.",
		RunE: func(_ *cobra.Command, args []string) error {
			o.Entries = args
			err := o.Run()
			if err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringSliceVarP(&o.Arguments, "argument", "D", []string{},
		i18n.T("Specify the top-level argument"))
	cmd.Flags().StringSliceVarP(&o.Settings, "setting", "Y", []string{},
		i18n.T("Specify the command line setting files"))
	cmd.Flags().StringVarP(&o.Output, "output", "o", "",
		i18n.T("Specify the output file"))
	cmd.Flags().BoolVarP(&o.DisableNone, "disable-none", "n", false,
		i18n.T("Disable dumping None values"))
	cmd.Flags().StringSliceVarP(&o.Overrides, "overrides", "O", []string{},
		i18n.T("Specify the configuration override path and value"))

	return cmd
}
