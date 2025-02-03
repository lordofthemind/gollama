package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lordofthemind/gollama/configs"
)

// DetermineModels selects the models to use based on flags
func DetermineModels(config configs.GollamaGlobalConfig, useAll, usePrimary, useSecondary, useTertiary bool, customModel string) ([]string, error) {
	if useAll {
		return []string{config.Primary.Model, config.Secondary.Model, config.Tertiary.Model}, nil
	}
	if customModel != "" {
		return []string{customModel}, nil
	}
	models := []string{}
	if usePrimary {
		models = append(models, config.Primary.Model)
	}
	if useSecondary {
		models = append(models, config.Secondary.Model)
	}
	if useTertiary {
		models = append(models, config.Tertiary.Model)
	}
	if len(models) == 0 {
		return nil, errors.New("no model specified")
	}
	return models, nil
}

// QueryAllModels queries all models and returns their responses
func QueryAllModels(models []string, prompt string, nonStreaming bool) map[string]string {
	responses := make(map[string]string)
	for _, model := range models {
		response := simulateChatResponse(model, prompt, nonStreaming)
		responses[model] = response
	}
	return responses
}

// PromptUserForModel lets the user select a model for further interaction
func PromptUserForModel(models []string) string {
	fmt.Println("Select a model for further conversation:")
	for i, model := range models {
		fmt.Printf("[%d] %s\n", i+1, model)
	}
	fmt.Print("Enter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	index := 0
	fmt.Sscanf(choice, "%d", &index)
	if index < 1 || index > len(models) {
		fmt.Println("Invalid choice.")
		return ""
	}
	return models[index-1]
}

// ChatWithModel facilitates a chat session with the selected model
func ChatWithModel(model, prompt string, nonStreaming bool) {
	response := simulateChatResponse(model, prompt, nonStreaming)
	fmt.Printf("[%s]: %s\n", model, response)
}

// Simulate chat response (placeholder for actual chat logic)
func simulateChatResponse(model, prompt string, nonStreaming bool) string {
	mode := "streaming"
	if nonStreaming {
		mode = "non-streaming"
	}
	return fmt.Sprintf("Response from %s in %s mode to prompt: %s", model, mode, prompt)
}
