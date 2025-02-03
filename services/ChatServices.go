package services

import (
	"errors"
	"fmt"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/helpers"
)

// InitiateChatSession starts a chat session based on user flags and input
func InitiateChatSession(
	config configs.GollamaGlobalConfig,
	configPath string,
	useAllModels, usePrimary, useSecondary, useTertiary bool,
	specificModel string,
	nonStreaming bool,
	prompt string,
) error {
	// Determine models to use
	models, err := helpers.DetermineModels(config, useAllModels, usePrimary, useSecondary, useTertiary, specificModel)
	if err != nil {
		return err
	}

	// Handle "use all models" mode
	if useAllModels {
		return handleAllModels(models, prompt, nonStreaming)
	}

	// Default to the primary model if no flags are set
	if len(models) == 0 {
		fmt.Println("No flags provided; using the primary model by default.")
		models = append(models, config.Primary.Model)
	}

	// Handle single or default model chat
	for _, model := range models {
		helpers.ChatWithModel(model, prompt, nonStreaming)
	}
	return nil
}

// handleAllModels manages the "use all models" mode
func handleAllModels(models []string, prompt string, nonStreaming bool) error {
	responses := helpers.QueryAllModels(models, prompt, nonStreaming)
	fmt.Println("Responses from all models:")
	for model, response := range responses {
		fmt.Printf("[%s]: %s\n", model, response)
	}

	selectedModel := helpers.PromptUserForModel(models)
	if selectedModel == "" {
		return errors.New("no model selected; aborting chat session")
	}

	fmt.Printf("Switching to model: %s\n", selectedModel)
	helpers.ChatWithModel(selectedModel, prompt, nonStreaming)
	return nil
}
