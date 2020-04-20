package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/v31/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type ghLabel struct {
	Labels []Label
}

type Label struct {
	Name        string
	Description string
	Color       string
}

var inputFile string

func init() {
	setCmd.PersistentFlags().StringVar(&inputFile, "label-file", "", "File that contains the labels to be applied in the GitHub org/repo.")

	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set GitHub Labels for a specific org/repo.",
	Long:  `Set GitHub Labels for a specific org/repo.`,
	Run: func(cmd *cobra.Command, args []string) {
		setGitHubLabels()
	},
}

func setGitHubLabels() {

	var t ghLabel

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Errorf("failed to read the file: %v", err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Errorf("failed to decode the file: %v", err.Error())
		os.Exit(1)
	}

	ghClient := gitHubClient()

	allLabels, err := listGHLabels(ghClient)
	if err != nil {
		log.Errorf("failed to get the existing labels: %v", err.Error())
		os.Exit(1)
	}

	for _, label := range t.Labels {

		log.Printf("Setting Github Label -> Name: %s\tDesc: %s\tColor:%s\n", label.Name, label.Description, label.Color)

		ghLabel := &github.Label{
			Name:        &label.Name,
			Description: &label.Description,
			Color:       &label.Color,
		}

		found := false
		for _, existingLabel := range allLabels {
			if strings.EqualFold(existingLabel.GetName(), label.Name) {
				if !isLabelEqual(label, existingLabel) {
					_, _, err = ghClient.Issues.EditLabel(context.Background(), org, repo, label.Name, ghLabel)
					if err != nil {
						log.Errorf("Failed to edit label: %v\n", err.Error())
						continue
					}

					log.Println("Edited label")
				} else {
					log.Println("Label is in sync")
				}

				found = true
			}
		}

		if !found {
			_, _, err := ghClient.Issues.CreateLabel(context.Background(), org, repo, ghLabel)
			if err != nil {
				log.Errorf("Failed to create label: %v\n", err.Error())
				continue
			}

			log.Println("Created label")
		}
	}

}

func isLabelEqual(label Label, gh *github.Label) bool {
	return gh.GetName() == label.Name &&
		gh.GetDescription() == label.Description &&
		gh.GetColor() == label.Color
}
