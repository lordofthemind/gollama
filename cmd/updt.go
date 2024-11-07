package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

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
			updateAllModels(models)
			return
		}

		// Handle the --model flag
		if modelFlag != "" {
			updateSpecificModel(modelFlag)
			return
		}

		// If no flag is provided, prompt the user for models to update
		promptUserForModels(models)
	},
}

func init() {
	rootCmd.AddCommand(updtCmd)

	updtCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Update all models")
	updtCmd.Flags().StringVarP(&modelFlag, "model", "m", "", "Specify a model to update")
}

// updateAllModels updates all available models by calling "ollama pull <model name>" for each model
func updateAllModels(models []string) {
	fmt.Println("Updating all models...")
	for _, model := range models {
		err := executeOllamaPull(model)
		if err != nil {
			fmt.Printf("Failed to update model %s: %v\n", model, err)
		} else {
			fmt.Printf("Successfully updated model %s.\n", model)
		}
	}
}

// updateSpecificModel updates a specific model provided by the user
func updateSpecificModel(model string) {
	fmt.Printf("Updating model %s...\n", model)
	err := executeOllamaPull(model)
	if err != nil {
		fmt.Printf("Failed to update model %s: %v\n", model, err)
	} else {
		fmt.Printf("Successfully updated model %s.\n", model)
	}
}

// promptUserForModels lists available models and allows the user to select models to update
func promptUserForModels(models []string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Available models:")
		for i, model := range models {
			fmt.Printf("%d. %s\n", i+1, model)
		}
		fmt.Print("Enter the number of the model to update, or type 'exit' to stop: ")

		// Read user input
		input, _ := reader.ReadString('\n')
		input = utils.TrimInput(input)

		// Handle exit condition
		if input == "exit" {
			fmt.Println("Exiting update prompt.")
			break
		}

		// Attempt to convert input to an integer for model selection
		index, err := utils.ParseInputIndex(input, len(models))
		if err != nil {
			fmt.Println("Invalid selection, please try again.")
			continue
		}

		// Update the selected model
		selectedModel := models[index]
		updateSpecificModel(selectedModel)
	}
}

// executeOllamaPull runs the "ollama pull <model name>" command
func executeOllamaPull(model string) error {
	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
