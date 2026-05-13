package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yoshi-komoto/gitbucket-cli/internal/output"
)

var prCommentListCmd = &cobra.Command{
	Use:   "list <number>",
	Short: "List comments on a pull request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil || n <= 0 {
			return fmt.Errorf("number must be a positive integer")
		}
		s, err := newSession(cmd.Context(), flags.repo)
		if err != nil {
			return err
		}
		cs, err := s.client.ListIssueComments(cmd.Context(), s.owner, s.repo, n)
		if err != nil {
			return err
		}
		return output.RenderCommentList(os.Stdout, flags.output, cs)
	},
}

func init() {
	prCommentCmd.AddCommand(prCommentListCmd)
}
