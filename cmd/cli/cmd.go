package main

import (
	"strings"

	"github.com/heathcliff26/netrouse/pkg/server"
	"github.com/heathcliff26/netrouse/pkg/version"
	"github.com/heathcliff26/netrouse/pkg/wol"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	name := strings.ToLower(version.Name)
	cobra.AddTemplateFunc(
		"ProgramName", func() string {
			return name
		},
	)

	rootCmd := &cobra.Command{
		Use:   name,
		Short: version.Name + " power on other devices on the network via Wake-on-Lan",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	wolCMD := wol.NewCommand()

	rootCmd.AddCommand(
		wolCMD,
		server.NewCommand(),
		version.NewCommand(),
	)

	return rootCmd
}
