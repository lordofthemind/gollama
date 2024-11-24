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
		return nil, fmt.Errorf("failed to execute 'ollama list': %w", err)
	}

	lines := strings.Split(string(output), "\n")
	models := []string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] != "NAME" {
			models = append(models, fields[0])
		}
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

// ValidateTemperature validates if the given temperature is within the allowed range
func ValidateTemperature(temp float64) error {
	if temp < 0.1 || temp > 1.0 {
		return fmt.Errorf("temperature must be between 0.1 and 1.0")
	}
	return nil
}
