package gb

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yoshi-komoto/gitbucket-cli/internal/output"
)

var prViewCmd = &cobra.Command{
	Use:   "view <number>",
	Short: "View a pull request",
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
		pr, err := s.client.GetPull(cmd.Context(), s.owner, s.repo, n)
		if err != nil {
			return err
		}
		return output.RenderPullView(os.Stdout, flags.output, pr)
	},
}

func init() {
	prCmd.AddCommand(prViewCmd)
}
