package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat with the assistant using the selected model",
	Long: `The chat command enables an interactive conversation with the assistant.
By default, it uses streaming mode with the primary model from the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the configuration
		config, err := configs.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Determine the model(s) to use based on flags
		useAllModels, _ := cmd.Flags().GetBool("all")
		models := []string{}
		var selectedModel string

		if useAllModels {
			// Use all models if -a/--all flag is set
			models = append(models, config.PrimaryModel, config.SecondaryModel, config.TertiaryModel)
			fmt.Println("Using all models. After the first response, you’ll be prompted to select a model to continue with.")
		} else {
			// Select single model based on flag precedence
			if useSecondary, _ := cmd.Flags().GetBool("secondary"); useSecondary {
				selectedModel = config.SecondaryModel
			} else if useTertiary, _ := cmd.Flags().GetBool("tertiary"); useTertiary {
				selectedModel = config.TertiaryModel
			} else {
				selectedModel = config.PrimaryModel
			}
			models = append(models, selectedModel)
		}

		// Set initial prompt from args or go into interactive mode if no prompt provided
		prompt := strings.Join(args, " ")
		if prompt == "" {
			fmt.Println("Interactive mode activated. Start chatting below:")
		}

		// Begin interactive chat session
		chatLoop(cmd, models, selectedModel, prompt)
	},
}

func chatLoop(cmd *cobra.Command, models []string, selectedModel, prompt string) {
	// Initialize reader for interactive input
	reader := bufio.NewReader(os.Stdin)

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

		// Handle model selection if using all models
		if len(models) > 1 && selectedModel == "" {
			fmt.Println("Available models:")
			for i, model := range models {
				fmt.Printf("[%d] %s\n", i+1, model)
			}
			fmt.Print("Select model to continue with (e.g., 1, 2, 3): ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			// Validate choice and set selectedModel
			idx := -1
			if n, err := fmt.Sscanf(choice, "%d", &idx); err == nil && n == 1 && idx >= 1 && idx <= len(models) {
				selectedModel = models[idx-1]
			} else {
				fmt.Println("Invalid choice. Please select a valid model number.")
				continue
			}
		}

		// Create a new context and API client
		ctx := context.Background()
		client, err := api.ClientFromEnvironment()
		if err != nil {
			log.Fatal(err)
		}

		// Set up the request with the selected model and prompt
		promptConfig := &api.GenerateRequest{
			Model:  selectedModel,
			Prompt: prompt,
		}

		// Check response mode and call appropriate function
		responseMode, _ := cmd.Flags().GetBool("response")
		if responseMode {
			// Non-streaming mode
			promptConfig.Stream = new(bool) // Non-streaming if false
			err = client.Generate(ctx, promptConfig, func(resp api.GenerateResponse) error {
				fmt.Println("Gollama:", resp.Response)
				return nil
			})
		} else {
			// Streaming mode
			var responseBuilder strings.Builder
			firstChunk := true
			err = client.Generate(ctx, promptConfig, func(resp api.GenerateResponse) error {
				// Only print "Gollama:" label once at the beginning
				if firstChunk {
					fmt.Print("Gollama: ")
					firstChunk = false
				}
				// Append each chunk of response text to the builder
				responseBuilder.WriteString(resp.Response)
				fmt.Print(resp.Response) // Print chunk directly
				return nil
			})
			fmt.Println() // Ensure a newline after the streaming response
		}

		if err != nil {
			fmt.Printf("Error generating response from model %s: %v\n", selectedModel, err)
		}

		// Reset prompt for next user input and continue the loop
		prompt = ""
	}
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Define flags for model selection and response mode
	chatCmd.Flags().BoolP("response", "r", false, "Use non-streaming mode")
	chatCmd.Flags().BoolP("secondary", "s", false, "Use secondary model")
	chatCmd.Flags().BoolP("tertiary", "t", false, "Use tertiary model")
	chatCmd.Flags().BoolP("all", "a", false, "Use all models")
}
