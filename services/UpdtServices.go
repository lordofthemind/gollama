package services

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lordofthemind/gollama/helpers"
)

// GetAvailableModels retrieves the list of models from the Ollama installation.
func GetAvailableModels() ([]string, error) {
	return helpers.GetOllamaModels()
}

// UpdateAllModels updates all models in the installation.
func UpdateAllModels(models []string) error {
	fmt.Println("Updating all models...")
	for _, model := range models {
		fmt.Printf("Trying to update %s!\n", model)
		if err := helpers.PullModel(model); err != nil {
			fmt.Printf("Failed to update model %s: %v\n", model, err)
		} else {
			fmt.Printf("Successfully updated model %s.\n", model)
		}
	}
	return nil
}

// UpdateSpecificModel updates a single specified model.
func UpdateSpecificModel(model string) error {
	fmt.Printf("Updating model %s...\n", model)
	return helpers.PullModel(model)
}

// PromptUserForModelSelection displays a list of models and allows the user to select models to update interactively.
func PromptUserForModelSelection(models []string) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Available models:")
		for i, model := range models {
			fmt.Printf("%d. %s\n", i+1, model)
		}
		fmt.Print("Enter the number of the model to update, or type 'exit' to stop: ")

		input, _ := reader.ReadString('\n')
		input = helpers.TrimWhitespace(input)

		if input == "exit" {
			fmt.Println("Exiting update prompt.")
			return nil
		}

		index, err := helpers.ParseModelIndex(input, len(models))
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			continue
		}

		selectedModel := models[index]
		if err := UpdateSpecificModel(selectedModel); err != nil {
			fmt.Printf("Failed to update model %s: %v\n", selectedModel, err)
		}
	}
}
