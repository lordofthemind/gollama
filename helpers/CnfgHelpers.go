package helpers

import (
	"fmt"
	"os/exec"
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
	return models, nil
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
