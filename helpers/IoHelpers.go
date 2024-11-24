package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NewReader creates a new reader for user input
func NewReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

// PromptForTemperature prompts the user for a temperature with default value
func PromptForTemperature(prompt string) float64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " (default: 0.4): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Use default temperature if input is empty
	if input == "" {
		return 0.4
	}

	temp, err := strconv.ParseFloat(input, 64)
	if err != nil || temp < 0.1 || temp > 1.0 {
		fmt.Println("Invalid temperature. It must be a number between 0.1 and 1.0. Please try again.")
		return PromptForTemperature(prompt)
	}

	return temp
}

// ReadInput reads input from the user
func ReadInput(reader *bufio.Reader) string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ConfirmAction asks the user to confirm an action with a "y" or "n" input
func ConfirmAction(prompt string) bool {
	for {
		fmt.Print(prompt)
		input := ReadInput(nil)
		input = strings.ToLower(input)
		if input == "y" {
			return true
		} else if input == "n" {
			return false
		}
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
	}
}

// ReadFloatInput validates and reads a float input from the user
func ReadFloatInput(reader *bufio.Reader, prompt string, min, max float64) float64 {
	for {
		fmt.Print(prompt)
		input := ReadInput(reader)
		if value, err := strconv.ParseFloat(input, 64); err == nil && value >= min && value <= max {
			return value
		}
		fmt.Printf("Please enter a valid number between %.2f and %.2f.\n", min, max)
	}
}
