package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/v31/github"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var saveFile bool

func init() {
	listCmd.PersistentFlags().BoolVar(&saveFile, "save", false, "Save the labels in a Yaml file format.")

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List GitHub Labels for a specific org/repo.",
	Long:  `List GitHub Labels for a specific org/repo.`,
	Run: func(cmd *cobra.Command, args []string) {
		listGitHubLabels()
	},
}

func listGitHubLabels() {
	ghClient := gitHubClient()

	allLabels, err := listGHLabels(ghClient)
	if err != nil {
		log.Errorf("failed to list labels: %v", err.Error())
		os.Exit(1)
	}

	var toFile ghLabel

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Label Name", "Description", "Color"})
	table.SetAutoWrapText(false)
	for _, label := range allLabels {
		table.Append([]string{strings.TrimSpace(label.GetName()), strings.TrimSpace(label.GetDescription()), strings.TrimSpace(label.GetColor())})

		toFile.Labels = append(toFile.Labels, Label{
			Name:        label.GetName(),
			Description: label.GetDescription(),
			Color:       label.GetColor(),
		})
	}
	table.Render()

	if saveFile {
		d, err := yaml.Marshal(&toFile)
		if err != nil {
			log.Errorf("failed to encode labels: %v\n", err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(fmt.Sprintf("./labels-%s-%s.yaml", org, repo), d, 0644)
		if err != nil {
			log.Errorf("failed to save label file err: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func listGHLabels(client *github.Client) ([]*github.Label, error) {
	var allLabels []*github.Label

	opt := &github.ListOptions{
		PerPage: 50,
	}

	for {
		repos, resp, err := client.Issues.ListLabels(context.Background(), org, repo, opt)
		if err != nil {
			log.Errorf("failed to list labels for %s/%s err: %s\n", org, repo, err.Error())
			return nil, err
		}

		allLabels = append(allLabels, repos...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allLabels, nil
}

func gitHubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitHubToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
