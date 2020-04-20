package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	gitHubToken string
	org         string
	repo        string

	rootCmd = &cobra.Command{
		Use:   "github-gitlab-labels",
		Short: "Tool to list/set github/gitlab labels for a specific org/repo.",
		Long:  `Tool to list/set github/gitlab labels for a specific org/repo.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(checkGHToken)

	rootCmd.PersistentFlags().StringVar(&gitHubToken, "github-token", "", "GitHub Token to be used to access the repository.")
	rootCmd.PersistentFlags().StringVar(&org, "org", "", "GitHub Org.")
	rootCmd.PersistentFlags().StringVar(&repo, "repo", "", "GitHub Repo.")

}

func checkGHToken() {
	if gitHubToken == "" {
		fmt.Println("Missing GitHub Token")
		os.Exit(1)
	}

	if org == "" {
		fmt.Println("Missing GitHub Org")
		os.Exit(1)
	}

	if repo == "" {
		fmt.Println("Missing GitHub Repo")
		os.Exit(1)
	}
}
