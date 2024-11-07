package utils

import (
	"fmt"
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
