package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/utils"
)

var (
	allFlag   bool
	modelFlag string
)

// updtCmd represents the updt command
var updtCmd = &cobra.Command{
	Use:   "updt",
	Short: "Update models in the Ollama installation",
	Long: `The updt command allows you to update one or more models in the Ollama installation.
Use the --all flag to update all models, or the --model flag to update a specific model.
If no flag is provided, you will be prompted to select models to update.`,
	Run: func(cmd *cobra.Command, args []string) {
		models, err := utils.GetOllamaModels()
		if err != nil {
			fmt.Println("Error retrieving models from Ollama:", err)
			return
		}

		// Handle the --all flag
		if allFlag {
			utils.UpdateAllModels(models)
			return
		}

		// Handle the --model flag
		if modelFlag != "" {
			utils.UpdateSpecificModel(modelFlag)
			return
		}

		// If no flag is provided, prompt the user for models to update
		utils.PromptUserForModels(models)
	},
}

func init() {
	rootCmd.AddCommand(updtCmd)

	updtCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Update all models")
	updtCmd.Flags().StringVarP(&modelFlag, "model", "m", "", "Specify a model to update")
}
