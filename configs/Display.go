package configs

import "fmt"

// DisplayGlobalConfig prints the current configuration to the console
func DisplayGlobalConfig(config GollamaGlobalConfig) {
	fmt.Println("Current Global Configuration:")
	fmt.Println("Primary:")
	fmt.Printf("  Model: %s\n", config.Primary.Model)
	fmt.Printf("  Temperature: %.2f\n", config.Primary.Temp)
	fmt.Println("Secondary:")
	fmt.Printf("  Model: %s\n", config.Secondary.Model)
	fmt.Printf("  Temperature: %.2f\n", config.Secondary.Temp)
	fmt.Println("Tertiary:")
	fmt.Printf("  Model: %s\n", config.Tertiary.Model)
	fmt.Printf("  Temperature: %.2f\n", config.Tertiary.Temp)
	fmt.Printf("Setup Completed: %t\n", config.SetupCompleted)
}
