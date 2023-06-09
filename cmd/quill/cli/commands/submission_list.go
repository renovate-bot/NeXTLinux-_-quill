package commands

import (
	"context"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/nextlinux/quill/cmd/quill/cli/application"
	"github.com/nextlinux/quill/cmd/quill/cli/options"
	"github.com/nextlinux/quill/internal/bus"
	"github.com/nextlinux/quill/internal/log"
	"github.com/nextlinux/quill/quill"
	"github.com/nextlinux/quill/quill/notary"
)

var _ options.Interface = &submissionListConfig{}

type submissionListConfig struct {
	options.Notary `yaml:"notary" json:"notary" mapstructure:"notary"`
}

func SubmissionList(app *application.Application) *cobra.Command {
	opts := &submissionListConfig{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list previous submissions to Apple's Notary service",
		Args:    cobra.NoArgs,
		PreRunE: app.Setup(opts),
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(cmd.Context(), async(func() error {
				log.Info("fetching previous submissions")

				cfg := quill.NewNotarizeConfig(
					opts.Notary.Issuer,
					opts.Notary.PrivateKeyID,
					opts.Notary.PrivateKey,
				)

				token, err := notary.NewSignedToken(cfg.TokenConfig)
				if err != nil {
					return err
				}

				a := notary.NewAPIClient(token, cfg.HTTPTimeout)

				sub := notary.ExistingSubmission(a, "")

				submissions, err := sub.List(context.Background())
				if err != nil {
					return err
				}

				// show list report

				t := table.NewWriter()
				t.SetStyle(table.StyleLight)

				t.AppendHeader(table.Row{"ID", "Name", "Status", "Created"})

				for _, item := range submissions {
					t.AppendRow(table.Row{item.ID, item.Name, item.Status, item.CreatedDate})
				}

				bus.Report(t.Render())

				return nil
			}))
		},
	}

	commonConfiguration(app, cmd, opts)

	return cmd
}
