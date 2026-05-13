package gb

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/yoshi-komoto/gitbucket-cli/internal/output"
)

var prCommentAddFlags struct {
	body     string
	bodyFile string
}

var prCommentAddCmd = &cobra.Command{
	Use:   "add <number>",
	Short: "Add a comment to a pull request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil || n <= 0 {
			return fmt.Errorf("number must be a positive integer")
		}
		body, err := resolveCommentBody(prCommentAddFlags.body, prCommentAddFlags.bodyFile, os.Stdin)
		if err != nil {
			return err
		}
		s, err := newSession(cmd.Context(), flags.repo)
		if err != nil {
			return err
		}
		created, err := s.client.CreateIssueComment(cmd.Context(), s.owner, s.repo, n, body)
		if err != nil {
			return err
		}
		return output.RenderComment(os.Stdout, flags.output, created)
	},
}

// errEmptyBody indicates the resolved body is empty.
var errEmptyBody = errors.New("empty body")

func resolveCommentBody(bodyFlag, bodyFile string, stdin *os.File) (string, error) {
	if bodyFlag != "" && bodyFile != "" {
		return "", fmt.Errorf("--body and --body-file are mutually exclusive")
	}
	switch {
	case bodyFlag != "":
		return validateBody(bodyFlag)
	case bodyFile != "":
		b, err := os.ReadFile(bodyFile)
		if err != nil {
			return "", fmt.Errorf("read body file: %w", err)
		}
		return validateBody(string(b))
	default:
		if stdin == nil || isatty.IsTerminal(stdin.Fd()) {
			return "", fmt.Errorf("no body provided. Pass --body, --body-file, or pipe via stdin")
		}
		b, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return validateBody(string(b))
	}
}

func validateBody(s string) (string, error) {
	s = strings.TrimRight(s, "\r\n\t ")
	if strings.TrimSpace(s) == "" {
		return "", errEmptyBody
	}
	return s, nil
}

func init() {
	prCommentAddCmd.Flags().StringVar(&prCommentAddFlags.body, "body", "", "Comment body")
	prCommentAddCmd.Flags().StringVar(&prCommentAddFlags.bodyFile, "body-file", "", "Path to a file containing the comment body")
	prCommentCmd.AddCommand(prCommentAddCmd)
}
