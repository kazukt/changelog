package cmd

import (
	"bytes"
	"os"

	"github.com/kazukt/changelog/pkg/changelog"
	"github.com/spf13/cobra"
)

func NewCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes a new changelog",
		Long: `Outputs an empty changelog, with preamble and Unreleased version
You can specify a filename using the --output/-o flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()

			c := changelog.Changelog{
				Title: "Changelog",
				Preamble: `All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).`,
				Versions: []*changelog.Version{
					{
						Name: "Unreleased",
						Changes: []*changelog.ChangeCollection{
							{
								Type: changelog.ChangeTypeAdded,
								Items: []string{
									"First commit.",
								},
							},
						},
					},
				},
			}

			var buf bytes.Buffer
			err := c.Write(&buf)
			if err != nil {
				return err
			}

			dest, err := flags.GetString("output")
			if err != nil {
				return err
			}
			if dest == "" {
				if _, err := buf.WriteTo(cmd.OutOrStdout()); err != nil {
					return err
				}
				return nil
			}
			f, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := buf.WriteTo(f); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd

}
