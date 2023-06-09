package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/nextlinux/quill/cmd/quill/cli/application"
	"github.com/nextlinux/quill/cmd/quill/cli/options"
	"github.com/nextlinux/quill/internal/bus"
	"github.com/nextlinux/quill/quill/extract"
)

var _ options.Interface = &describeConfig{}

type describeConfig struct {
	Path             string `yaml:"path" json:"path" mapstructure:"path"`
	options.Format   `yaml:",inline" json:",inline" mapstructure:",squash"`
	options.Describe `yaml:"describe" json:"describe" mapstructure:"describe"`
}

func (d *describeConfig) Redact() {
	options.RedactAll(&d.Format, &d.Describe)
}

func (d *describeConfig) AddFlags(flags *pflag.FlagSet) {
	options.AddAllFlags(flags, &d.Format, &d.Describe)
}

func (d *describeConfig) BindFlags(flags *pflag.FlagSet, v *viper.Viper) error {
	return options.BindAllFlags(flags, v, &d.Format, &d.Describe)
}

func Describe(app *application.Application) *cobra.Command {
	opts := &describeConfig{
		Format: options.Format{
			Output:           "text",
			AllowableFormats: []string{"text", "json"},
		},
	}

	cmd := &cobra.Command{
		Use:   "describe PATH",
		Short: "show the details of a macho binary",
		Example: options.FormatPositionalArgsHelp(
			map[string]string{
				"PATH": "the darwin binary to print details for",
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
				var err error
				buf := &strings.Builder{}
				switch strings.ToLower(opts.Output) {
				case "text":
					err = extract.ShowText(opts.Path, buf, !opts.Detail)
				case "json":
					err = extract.ShowJSON(opts.Path, buf)
				default:
					err = fmt.Errorf("unknown format: %s", opts.Output)
				}

				if err != nil {
					return err
				}

				bus.Report(buf.String())

				return nil
			}))
		},
	}

	commonConfiguration(app, cmd, opts)

	return cmd
}
