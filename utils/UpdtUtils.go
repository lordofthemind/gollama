package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

)

// ParseInputIndex converts a user input string to a valid index within the bounds of a list
func ParseInputIndex(input string, maxIndex int) (int, error) {
	// Convert the input to an integer
	index, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("invalid input: please enter a number")
	}

	// Check if the index is within valid range
	if index < 1 || index > maxIndex {
		return -1, fmt.Errorf("out of range: please enter a number between 1 and %d", maxIndex)
	}

	// Return the zero-based index
	return index - 1, nil
}

// TrimInput reads a string input and trims any leading or trailing whitespace
func TrimInput(input string) string {
	return strings.TrimSpace(input)
}

// updateAllModels updates all available models by calling "ollama pull <model name>" for each model
func UpdateAllModels(models []string) {
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
func UpdateSpecificModel(model string) {
	fmt.Printf("Updating model %s...\n", model)
	err := executeOllamaPull(model)
	if err != nil {
		fmt.Printf("Failed to update model %s: %v\n", model, err)
	} else {
		fmt.Printf("Successfully updated model %s.\n", model)
	}
}

// promptUserForModels lists available models and allows the user to select models to update
func PromptUserForModels(models []string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Available models:")
		for i, model := range models {
			fmt.Printf("%d. %s\n", i+1, model)
		}
		fmt.Print("Enter the number of the model to update, or type 'exit' to stop: ")

		// Read user input
		input, _ := reader.ReadString('\n')
		input = TrimInput(input)

		// Handle exit condition
		if input == "exit" {
			fmt.Println("Exiting update prompt.")
			break
		}

		// Attempt to convert input to an integer for model selection
		index, err := ParseInputIndex(input, len(models))
		if err != nil {
			fmt.Println("Invalid selection, please try again.")
			continue
		}

		// Update the selected model
		selectedModel := models[index]
		UpdateSpecificModel(selectedModel)
	}
}

// executeOllamaPull runs the "ollama pull <model name>" command
func executeOllamaPull(model string) error {
	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
