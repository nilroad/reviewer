package command

import (
	"club/internal/version"
	"fmt"

	"github.com/spf13/cobra"
)

type Version struct{}

func (cmd *Version) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version and commit hash",
		Long:  `print version and commit hash`,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(
				"Version: %v\nRelease Date: %v\nCommit Hash: %v\n",
				version.Version,
				version.ReleaseDate,
				version.CommitHash,
			)
		},
	}
}
