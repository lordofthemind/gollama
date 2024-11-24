package helpers

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetOllamaModels retrieves the available models from the Ollama installation
func GetOllamaModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	models := []string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] != "NAME" {
			models = append(models, fields[0])
		}
	}

	// Check if no models were found
	if len(models) == 0 {
		fmt.Println("No models found. You can pull models using the command:")
		fmt.Println("  ollama pull <model_name>")
	}

	return models, nil
}

// SelectModel allows user to select a model by index
func SelectModel(reader *bufio.Reader, models []string, prompt string) string {
	for {
		fmt.Println(prompt)
		for i, model := range models {
			fmt.Printf("%d. %s\n", i+1, model)
		}
		fmt.Print("Select a model by entering the number: ")
		input := ReadInput(reader)
		modelIndex, err := strconv.Atoi(input)
		if err == nil && modelIndex > 0 && modelIndex <= len(models) {
			return models[modelIndex-1]
		}
		fmt.Println("Invalid choice. Please select a valid number.")
	}
}

// ValidateModel checks if the model exists in the list of available models
func ValidateModel(model string, models []string) bool {
	for _, m := range models {
		if m == model {
			return true
		}
	}
	fmt.Printf("Model '%s' is not available.\n", model)
	return false
}
