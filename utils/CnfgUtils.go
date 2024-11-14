package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Function to retrieve models available in `ollama`
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
		if len(fields) > 0 && fields[0] != "NAME" { // Ignore header line
			models = append(models, fields[0])
		}
	}
	return models, nil
}

// Helper function to display temperature guidance
func DisplayTemperatureGuidance() {
	fmt.Println("\n### Temperature Guidance ###")
	fmt.Println("Temperature settings guide:")
	fmt.Println("0.0 - 0.3: Deterministic, ideal for precise tasks")
	fmt.Println("0.4 - 0.7: Balanced, suitable for conversations")
	fmt.Println("0.8 - 1.0: High randomness, for creative tasks")
}

// Helper function to select a model based on user input
func SelectModel(reader *bufio.Reader, models []string) string {
	for {
		modelIndex, _ := strconv.Atoi(ReadInput(reader))
		if modelIndex > 0 && modelIndex <= len(models) {
			return models[modelIndex-1]
		} else {
			fmt.Printf("Invalid choice. Please select a number between 1 and %d: ", len(models))
		}
	}
}

// readInput reads input from the user and trims whitespace
func ReadInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// readFloatInput reads a float input from the user
func ReadFloatInput(reader *bufio.Reader) float64 {
	var value float64
	for {
		fmt.Print("Enter Temperature (0.1 - 1.0, e.g., 0.5): ")
		input := ReadInput(reader)
		_, err := fmt.Sscanf(input, "%f", &value)
		if err == nil && value >= 0.1 && value <= 1.0 {
			break
		}
		fmt.Println("Invalid input. Please enter a valid number between 0.1 and 1.0.")
	}
	return value
}

func ValidateModel(model string, models []string) bool {
	for _, m := range models {
		if m == model {
			return true
		}
	}
	fmt.Printf("Model %s is not available in Ollama installation.\n", model)
	return false
}
