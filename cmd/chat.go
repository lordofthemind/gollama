package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat with the assistant using the selected model",
	Long: `Start an interactive conversation with the assistant.
Flags allow specifying which model(s) to use for the chat.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve flags
		useAllModels, _ := cmd.Flags().GetBool("all")
		usePrimary, _ := cmd.Flags().GetBool("primary")
		useSecondary, _ := cmd.Flags().GetBool("secondary")
		useTertiary, _ := cmd.Flags().GetBool("tertiary")
		specificModel, _ := cmd.Flags().GetString("model")
		nonStreaming, _ := cmd.Flags().GetBool("response")

		// Combine the prompt from args
		prompt := strings.Join(args, " ")
		if prompt == "" {
			fmt.Println("Error: No prompt provided. Please provide a prompt to start the chat.")
			os.Exit(1)
		}

		// Load configuration
		config, configPath, err := configs.LoadGlobalConfig()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Invoke the service layer to handle chat interactions
		err = services.InitiateChatSession(config, configPath, useAllModels, usePrimary, useSecondary, useTertiary, specificModel, nonStreaming, prompt)
		if err != nil {
			fmt.Printf("Error during chat: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Define flags for chat options
	chatCmd.Flags().BoolP("response", "r", false, "Use non-streaming response mode")
	chatCmd.Flags().BoolP("all", "a", false, "Use all models for the prompt")
	chatCmd.Flags().BoolP("primary", "p", false, "Use the primary model")
	chatCmd.Flags().BoolP("secondary", "s", false, "Use the secondary model")
	chatCmd.Flags().BoolP("tertiary", "t", false, "Use the tertiary model")
	chatCmd.Flags().StringP("model", "m", "", "Specify a custom model")
}
