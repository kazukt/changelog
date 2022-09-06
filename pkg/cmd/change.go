package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kazukt/changelog/pkg/changelog"
	"github.com/kazukt/changelog/pkg/parser"
	"github.com/spf13/cobra"
)

func NewCmdChangeType(ctype changelog.ChangeType) *cobra.Command {
	name := string(ctype)

	var option struct {
		message string
	}

	cmd := &cobra.Command{
		Use:   strings.ToLower(name),
		Short: fmt.Sprintf(`Add item under %q change`, name),
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
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

			c.AddItem(ctype, option.message)

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
	flags.StringVarP(&option.message, "message", "m", "", "Change message")

	return cmd

}
