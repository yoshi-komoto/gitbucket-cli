package cli

import "github.com/spf13/cobra"

var prCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage pull request comments",
}

func init() {
	prCmd.AddCommand(prCommentCmd)
}
