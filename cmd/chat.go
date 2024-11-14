package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/utils"
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
		var modelLocked bool // Tracks if a specific model has been locked after the first response

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
			fmt.Printf("Interactive mode activated with model: %s. Start chatting below:\n", selectedModel)
		}

		// Begin interactive chat session
		utils.ChatLoop(cmd, models, &selectedModel, &modelLocked, prompt)
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
