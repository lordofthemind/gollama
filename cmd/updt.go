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

Key Features:
- Update all models using the '--all' flag for a complete refresh.
- Update a specific model using the '--model' flag with the model's name.
- Interactive selection mode if no flags are provided, enabling you to choose which models to update.
- Seamlessly integrates with Ollama's list and pull commands for efficient updates.

Usage:
- Run 'updt' without any flags to launch the interactive update mode.
- Use '--all' or '-a' to update all models in one go.
- Use '--model <model-name>' or '-m <model-name>' to update a specific model directly.`,

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
