package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/yoshi-komoto/gitbucket-cli/internal/gitbucket"
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
)

const bodyMaxColumn = 60

func RenderPullList(w io.Writer, format string, prs []gitbucket.PullRequest) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, prs)
	case FormatTable, "":
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "#\tTITLE\tSTATE\tAUTHOR\tUPDATED")
		for _, pr := range prs {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n",
				pr.Number,
				truncate(pr.Title, bodyMaxColumn),
				pr.State,
				pr.User.Login,
				pr.UpdatedAt.UTC().Format("2006-01-02"),
			)
		}
		return tw.Flush()
	default:
		return fmt.Errorf("unknown output format %q", format)
	}
}

func RenderPullView(w io.Writer, format string, pr *gitbucket.PullRequest) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, pr)
	case FormatTable, "":
		fmt.Fprintf(w, "#%d  %s\n", pr.Number, pr.Title)
		fmt.Fprintf(w, "state: %s    author: %s    base: %s ← %s\n",
			pr.State, pr.User.Login, pr.Base.Ref, pr.Head.Ref)
		fmt.Fprintf(w, "created: %s  updated: %s\n",
			pr.CreatedAt.UTC().Format("2006-01-02"),
			pr.UpdatedAt.UTC().Format("2006-01-02"))
		fmt.Fprintln(w)
		fmt.Fprintln(w, "--")
		fmt.Fprintln(w, pr.Body)
		return nil
	default:
		return fmt.Errorf("unknown output format %q", format)
	}
}

func RenderCommentList(w io.Writer, format string, cs []gitbucket.IssueComment) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, cs)
	case FormatTable, "":
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tAUTHOR\tCREATED\tBODY")
		for _, c := range cs {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n",
				c.ID,
				c.User.Login,
				c.CreatedAt.UTC().Format("2006-01-02 15:04"),
				truncate(flatten(c.Body), bodyMaxColumn),
			)
		}
		return tw.Flush()
	default:
		return fmt.Errorf("unknown output format %q", format)
	}
}

func RenderComment(w io.Writer, format string, c *gitbucket.IssueComment) error {
	return RenderCommentList(w, format, []gitbucket.IssueComment{*c})
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func flatten(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

func truncate(s string, max int) string {
	if max <= 0 {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-3]) + "..."
}
