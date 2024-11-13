package cmd

import (
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

		if useAllModels {
			// Use all models if -a/--all flag is set
			models = append(models, config.PrimaryModel, config.SecondaryModel, config.TertiaryModel)
		} else {
			// Select single model based on flag precedence
			if useSecondary, _ := cmd.Flags().GetBool("secondary"); useSecondary {
				models = append(models, config.SecondaryModel)
			} else if useTertiary, _ := cmd.Flags().GetBool("tertiary"); useTertiary {
				models = append(models, config.TertiaryModel)
			} else {
				models = append(models, config.PrimaryModel)
			}
		}

		// Ensure a prompt is provided
		if len(args) == 0 {
			fmt.Println("Please provide a prompt for the chat command.")
			return
		}
		prompt := strings.Join(args, " ")

		// Execute response generation for each selected model
		for _, model := range models {
			fmt.Printf("Response from model %s:\n", model)

			// Set up the context and prompt configuration
			ctx := context.Background()
			client, err := api.ClientFromEnvironment()
			if err != nil {
				log.Fatal(err)
			}
			promptConfig := &api.GenerateRequest{
				Model:  model,
				Prompt: prompt,
			}

			// Check response mode and call appropriate function
			responseMode, _ := cmd.Flags().GetBool("response")
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
			}

			if err != nil {
				fmt.Printf("Error generating response from model %s: %v\n", model, err)
			}

			// Add newline to separate responses if multiple models are used
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Define flags for model selection and response mode
	chatCmd.Flags().BoolP("response", "r", false, "Use non-streaming mode")
	chatCmd.Flags().BoolP("secondary", "s", false, "Use secondary model")
	chatCmd.Flags().BoolP("tertiary", "t", false, "Use tertiary model")
	chatCmd.Flags().BoolP("all", "a", false, "Use all models")
}
