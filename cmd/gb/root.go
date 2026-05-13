package gb

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "gb",
	Short:         "GitBucket CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and returns the appropriate exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
