package cmd

import (
	"io"
	"os"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/kazukt/changelog/pkg/parser"
	"github.com/spf13/cobra"
)

func NewCmdRelease() *cobra.Command {

	const dateFormat = "2006-01-02"

	var option struct {
		date string
	}

	cmd := &cobra.Command{
		Use:   "release [version]",
		Short: "Change Unreleased to [version]",
		Long: `Change Unreleased section to [version], updating the compare links accordingly.
It will normalize the output with the new version.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()

			date, err := time.Parse(dateFormat, option.date)
			if err != nil {
				return err
			}
			semVersion, err := semver.NewVersion(args[0])
			if err != nil {
				return err
			}

			filename, err := flags.GetString("filename")
			if err != nil {
				return err
			}
			var r io.Reader
			if filename == "" {
				r = cmd.InOrStdin()
			} else {
				f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				r = f
			}
			c, err := parser.Parse(r)
			if err != nil {
				return err
			}

			c.Release(semVersion.String(), date.Format(dateFormat))

			output, err := flags.GetString("output")
			if err != nil {
				return err
			}
			var w io.Writer
			if output == "" {
				w = cmd.OutOrStdout()
			} else {
				f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				w = f
			}
			if err := c.Write(w); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()

	today := time.Now().Format(dateFormat)
	flags.StringVarP(&option.date, "date", "d", today, "Release date")

	return cmd
}
