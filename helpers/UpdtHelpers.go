package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// PullModel runs the `ollama pull <model>` command to update a specific model.
func PullModel(model string) error {
	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ParseModelIndex validates the user's input index for selecting a model.
func ParseModelIndex(input string, maxIndex int) (int, error) {
	index, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("input is not a valid number")
	}
	if index < 1 || index > maxIndex {
		return -1, fmt.Errorf("input out of range: valid range is 1 to %d", maxIndex)
	}
	return index - 1, nil
}

// TrimWhitespace trims leading and trailing whitespace from a string.
func TrimWhitespace(input string) string {
	return strings.TrimSpace(input)
}
