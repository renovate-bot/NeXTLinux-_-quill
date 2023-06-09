package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/nextlinux/quill/cmd/quill/cli/application"
	"github.com/nextlinux/quill/cmd/quill/cli/options"
	"github.com/nextlinux/quill/internal/log"
)

var _ options.Interface = &signConfig{}

type signAndNotarizeConfig struct {
	Path            string `yaml:"path" json:"path" mapstructure:"path"`
	options.Signing `yaml:"sign" json:"sign" mapstructure:"sign"`
	options.Notary  `yaml:"notary" json:"notary" mapstructure:"notary"`
	options.Status  `yaml:"status" json:"status" mapstructure:"status"`
	DryRun          bool `yaml:"dry-run" json:"dry-run" mapstructure:"dry-run"`
}

func (o *signAndNotarizeConfig) Redact() {
	options.RedactAll(&o.Notary, &o.Status, &o.Signing)
}

func (o *signAndNotarizeConfig) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&o.DryRun, "dry-run", o.DryRun, "dry run mode (do not actually notarize)")
	options.AddAllFlags(flags, &o.Notary, &o.Status, &o.Signing)
}

func (o *signAndNotarizeConfig) BindFlags(flags *pflag.FlagSet, v *viper.Viper) error {
	if err := options.Bind(v, "dry-run", flags.Lookup("dry-run")); err != nil {
		return err
	}
	return options.BindAllFlags(flags, v, &o.Notary, &o.Status, &o.Signing)
}

func SignAndNotarize(app *application.Application) *cobra.Command {
	opts := &signAndNotarizeConfig{
		Status:  options.DefaultStatus(),
		Signing: options.DefaultSigning(),
	}

	cmd := &cobra.Command{
		Use:   "sign-and-notarize PATH",
		Short: "sign and notarize a macho (darwin) executable binary",
		Example: options.FormatPositionalArgsHelp(
			map[string]string{
				"PATH": "the darwin binary to sign and notarize",
			},
		),
		Args: chainArgs(
			cobra.ExactArgs(1),
			func(_ *cobra.Command, args []string) error {
				opts.Path = args[0]
				return nil
			},
		),
		PreRunE: app.Setup(opts),
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(cmd.Context(), async(func() error {
				err := sign(opts.Path, opts.Signing)
				if err != nil {
					return fmt.Errorf("signing failed: %w", err)
				}

				if opts.DryRun {
					log.Warn("[DRY RUN] skipping notarization...")
					return nil
				}

				_, err = notarize(opts.Path, opts.Notary, opts.Status)
				if err != nil {
					return fmt.Errorf("notarization failed: %w", err)
				}

				return nil
			}))
		},
	}

	commonConfiguration(app, cmd, opts)

	return cmd
}
