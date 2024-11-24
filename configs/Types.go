package configs

// Primary, Secondary, Tertiary configurations
type ModelConfig struct {
	Model string  `mapstructure:"model"`
	Temp  float64 `mapstructure:"temp"`
}

// GollamaGlobalConfig represents the main configuration structure
type GollamaGlobalConfig struct {
	Primary        ModelConfig `mapstructure:"primary"`
	Secondary      ModelConfig `mapstructure:"secondary"`
	Tertiary       ModelConfig `mapstructure:"tertiary"`
	SetupCompleted bool        `mapstructure:"setup_completed"`
}
