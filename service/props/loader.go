package props

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	v := viper.New()

	// APP_ENV=dev → loads application-dev.yaml
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "default"
	}

	fileName := "application.yaml"
	if env != "default" {
		fileName = fmt.Sprintf("app-%s", env)
	}

	// config directory
	v.AddConfigPath("../../env")
	v.SetConfigName(fileName)
	v.SetConfigType("yaml")

	// enable environment overrides
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// read YAML
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	return cfg, nil
}
