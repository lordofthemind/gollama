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

// ReadInput reads and trims user input
func ReadInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ConfirmAction asks for a yes/no confirmation from the user
func ConfirmAction(prompt string) bool {
	fmt.Print(prompt + " (y/n): ")
	return ReadInput(NewReader()) == "y"
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
