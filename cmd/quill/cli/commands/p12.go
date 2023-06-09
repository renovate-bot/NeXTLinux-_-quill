package commands

import (
	"github.com/spf13/cobra"

	"github.com/nextlinux/quill/cmd/quill/cli/application"
)

func P12(_ *application.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "p12",
		Short: "describe and manipulate p12 files",
		Args:  cobra.NoArgs,
	}

	commonConfiguration(nil, cmd, nil)
	return cmd
}
