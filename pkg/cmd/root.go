package cmd

import (
	"github.com/kazukt/changelog/pkg/changelog"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "changelog <command> [flags]",
		Short:         "changelog CLI",
		Long:          `Manipulate and validate a Markdown changelog file from the command line.`,
		SilenceErrors: true,
	}

	// Child command
	cmd.AddCommand(NewCmdInit())
	cmd.AddCommand(NewCmdChangeType(changelog.ChangeTypeAdded))
	cmd.AddCommand(NewCmdChangeType(changelog.ChangeTypeChanged))
	cmd.AddCommand(NewCmdChangeType(changelog.ChangeTypeDeprecated))
	cmd.AddCommand(NewCmdChangeType(changelog.ChangeTypeFixed))
	cmd.AddCommand(NewCmdChangeType(changelog.ChangeTypeRemoved))
	cmd.AddCommand(NewCmdChangeType(changelog.ChangeTypeSecurity))
	cmd.AddCommand(NewCmdRelease())

	flags := cmd.PersistentFlags()
	flags.StringP("filename", "f", "", "Changelog file or stdin")
	cmd.MarkFlagFilename("filename")
	flags.StringP("output", "o", "", "Output file or stdout")
	cmd.MarkFlagFilename("output")

	return cmd
}
