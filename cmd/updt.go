package cmd

import (
	"fmt"

	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

var (
	allFlag   bool
	modelFlag string
)

// updtCmd represents the update command
var updtCmd = &cobra.Command{
	Use:   "updt",
	Short: "Update models in the Ollama installation",
	Long: `The updt command allows you to update one or more models in the Ollama installation.
Use the --all flag to update all models, or the --model flag to update a specific model.
If no flag is provided, you will be prompted to select models to update.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := handleUpdateCommand()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updtCmd)

	updtCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Update all models")
	updtCmd.Flags().StringVarP(&modelFlag, "model", "m", "", "Specify a model to update")
}

func handleUpdateCommand() error {
	models, err := services.GetAvailableModels()
	if err != nil {
		return fmt.Errorf("failed to retrieve models: %w", err)
	}

	if allFlag {
		return services.UpdateAllModels(models)
	}

	if modelFlag != "" {
		return services.UpdateSpecificModel(modelFlag)
	}

	return services.PromptUserForModelSelection(models)
}
