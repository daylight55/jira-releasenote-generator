package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var (
    Version   = "dev"
    Commit    = "none"
    Date      = "unknown"
    GoVersion = "unknown"
)

func newVersionCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Print version information",
        Long:  "Print detailed version information about this build of jira-changelog-generator",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("jira-changelog-generator version %s\n", Version)
            fmt.Printf("  Built with: %s\n", GoVersion)
            fmt.Printf("  Commit: %s\n", Commit)
            fmt.Printf("  Built on: %s\n", Date)
        },
    }
}