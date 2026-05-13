package gb

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yoshi-komoto/gitbucket-cli/internal/gitbucket"
	"github.com/yoshi-komoto/gitbucket-cli/internal/output"
)

var prListFlags struct {
	state string
	limit int
}

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if prListFlags.limit <= 0 || prListFlags.limit > 100 {
			return fmt.Errorf("--limit must be between 1 and 100")
		}
		switch prListFlags.state {
		case "open", "closed", "all":
		default:
			return fmt.Errorf("--state must be open, closed, or all")
		}
		s, err := newSession(cmd.Context(), flags.repo)
		if err != nil {
			return err
		}
		prs, err := s.client.ListPulls(cmd.Context(), s.owner, s.repo, gitbucket.PullsListOptions{
			State:   prListFlags.state,
			PerPage: prListFlags.limit,
		})
		if err != nil {
			return err
		}
		return output.RenderPullList(os.Stdout, flags.output, prs)
	},
}

func init() {
	prListCmd.Flags().StringVar(&prListFlags.state, "state", "open", "Filter by state: open|closed|all")
	prListCmd.Flags().IntVar(&prListFlags.limit, "limit", 30, "Max results (1-100)")
	prCmd.AddCommand(prListCmd)
}
