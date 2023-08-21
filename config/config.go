package config

type Config struct {
	//LogLevel string `mapstructure:"log_level"`
}

func NewConfig() Config {
	return Config{}
}

func DefaultConfig() Config {
	cfg := NewConfig()

	return cfg
}
