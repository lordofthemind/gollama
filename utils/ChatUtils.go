package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

func ChatLoop(cmd *cobra.Command, models []string, selectedModel *string, modelLocked *bool, prompt string) {
	// Initialize reader for interactive input
	reader := bufio.NewReader(os.Stdin)
	firstResponse := true

	// Recursive interactive loop
	for {
		// If prompt is empty, ask for user input
		if prompt == "" {
			fmt.Print("You: ")
			input, _ := reader.ReadString('\n')
			prompt = strings.TrimSpace(input)
		}

		// Exit chat if user types "bye" or "exit"
		if strings.EqualFold(prompt, "bye") || strings.EqualFold(prompt, "exit") {
			fmt.Println("Goodbye!")
			break
		}

		// Display responses from all models initially and ask user to select one
		if len(models) > 1 && firstResponse {
			firstResponse = false                         // Set firstResponse to false after the initial response
			displayResponsesFromAllModels(models, prompt) // Show initial responses from all models

			// Prompt user to choose a model for further responses
			fmt.Println("Do you want to continue with a specific model, or receive responses from all models?")
			fmt.Println("[1] Continue with Primary Model")
			fmt.Println("[2] Continue with Secondary Model")
			fmt.Println("[3] Continue with Tertiary Model")
			fmt.Println("[A] Continue with All Models")

			fmt.Print("Choose an option (1, 2, 3, A): ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			// Handle user choice to lock a model or continue with all
			switch choice {
			case "1":
				*selectedModel = models[0]
				*modelLocked = true
			case "2":
				*selectedModel = models[1]
				*modelLocked = true
			case "3":
				*selectedModel = models[2]
				*modelLocked = true
			case "A", "a":
				*selectedModel = "" // Keep responses from all models
			default:
				fmt.Println("Invalid choice. Continuing with responses from all models.")
			}

			// After selecting a model, wait for the next user prompt
			prompt = ""
			continue
		}

		// Generate response based on model selection for subsequent prompts
		if *selectedModel != "" { // Specific model selected
			fmt.Printf("Gollama [%s]: ", *selectedModel) // Display model name
			displayResponseFromModel(*selectedModel, prompt, cmd)
		} else { // All models selected
			displayResponsesFromAllModels(models, prompt)
		}

		// Reset prompt for next user input and continue the loop
		prompt = ""
	}
}

func displayResponseFromModel(model, prompt string, cmd *cobra.Command) {
	ctx := context.Background()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Set up the request with the selected model and prompt
	promptConfig := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}

	// Check response mode only if cmd is not nil
	responseMode := false
	if cmd != nil {
		responseMode, _ = cmd.Flags().GetBool("response")
	}

	if responseMode {
		// Non-streaming mode
		promptConfig.Stream = new(bool) // Non-streaming if false
		err = client.Generate(ctx, promptConfig, func(resp api.GenerateResponse) error {
			fmt.Println(resp.Response)
			return nil
		})
	} else {
		// Streaming mode
		err = client.Generate(ctx, promptConfig, func(resp api.GenerateResponse) error {
			fmt.Print(resp.Response)
			return nil
		})
		fmt.Println()
	}

	if err != nil {
		fmt.Printf("Error generating response from model %s: %v\n", model, err)
	}
}

func displayResponsesFromAllModels(models []string, prompt string) {
	fmt.Println("Gollama Responses from All Models:")
	for _, model := range models {
		fmt.Printf("\n[%s]:\n", model)
		displayResponseFromModel(model, prompt, nil) // Pass nil for cmd here
	}
	fmt.Println()
}
