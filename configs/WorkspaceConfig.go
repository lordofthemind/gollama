package configs

// GollamaWorkspaceConfig represents the structure of the configuration file
type GollamaWorkspaceConfig struct {
	ProjectName        string  `mapstructure:"project_name"`
	DefaultModel       string  `mapstructure:"default_model"`
	Temperature        float64 `mapstructure:"temperature"`
	SecondDefaultModel string  `mapstructure:"second_default_model"`
	ThirdDefaultModel  string  `mapstructure:"third_default_model"`
}
