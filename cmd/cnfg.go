package cmd

import (
	"fmt"

	"github.com/lordofthemind/gollama/helpers"
	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

var (
	tempFlag   float64
	pModelFlag string
	sModelFlag string
	tModelFlag string
)

// cnfgCmd represents the configuration management command
var cnfgCmd = &cobra.Command{
	Use:   "cnfg",
	Short: "Manage the application configuration",
	Long:  "The cnfg command allows you to manage Gollama's configuration, including models and temperature settings.",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := services.LoadConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		if !config.SetupCompleted {
			services.SetupConfiguration(&config, configPath)
			return
		}

		if !cmd.Flags().Changed("temp") && !cmd.Flags().Changed("pmodel") &&
			!cmd.Flags().Changed("smodel") && !cmd.Flags().Changed("tmodel") {
			fmt.Println("Current Configuration:")
			services.DisplayConfig(config)

			if helpers.ConfirmAction("Do you want to edit the configuration?") {
				services.SetupConfiguration(&config, configPath)
			} else {
				fmt.Println("No changes made.")
			}
			return
		}

		if services.UpdateConfigFromFlags(&config, tempFlag, pModelFlag, sModelFlag, tModelFlag) {
			if err := services.SaveConfig(config, configPath); err != nil {
				fmt.Println("Error saving configuration:", err)
			} else {
				fmt.Println("Configuration updated successfully.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cnfgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cnfgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cnfgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cnfgCmd.Flags().Float64VarP(&tempFlag, "temp", "t", 0.5, "Set the model temperature (0.1-1.0)")
	cnfgCmd.Flags().StringVarP(&pModelFlag, "primary", "p", "", "Set the primary model")
	cnfgCmd.Flags().StringVarP(&sModelFlag, "secondary", "s", "", "Set the secondary model")
	cnfgCmd.Flags().StringVarP(&tModelFlag, "tertiary", "e", "", "Set the tertiary model")
}
