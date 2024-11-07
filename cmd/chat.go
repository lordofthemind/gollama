package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/services"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat with the assistant using the selected model",
	Long: `The chat command enables an interactive conversation with the assistant.
By default, it uses streaming mode with the primary model from the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the configuration with workspace taking priority if available
		config, err := configs.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Determine the model and mode based on flags
		var model string
		if useSecondary, _ := cmd.Flags().GetBool("secondary"); useSecondary {
			model = config.SecondaryModel
		} else if useTertiary, _ := cmd.Flags().GetBool("tertiary"); useTertiary {
			model = config.TertiaryModel
		} else {
			model = config.PrimaryModel
		}

		// Check if streaming or non-streaming mode is requested
		if responseMode, _ := cmd.Flags().GetBool("response"); responseMode {
			// Non-streaming mode
			err := services.GenerateCompletion(model, args[0], config.Temperature)
			if err != nil {
				fmt.Printf("Error in completion generation: %v\n", err)
			}
		} else {
			// Default to streaming mode
			err := services.GenerateStreamingCompletion(model, args[0], config.Temperature)
			if err != nil {
				fmt.Printf("Error in streaming completion generation: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Define flags for model selection and response mode
	chatCmd.Flags().BoolP("response", "r", false, "Use non-streaming mode")
	chatCmd.Flags().BoolP("secondary", "s", false, "Use secondary model")
	chatCmd.Flags().BoolP("tertiary", "t", false, "Use tertiary model")
}
