package config

import (
	"fmt"
	"github.com/spf13/viper"

	"os"
)

type (
	Config struct {
		App      `mapstructure:"app"`
		HTTP     `mapstructure:"http"`
		Services `mapstructure:"services"`
	}

	App struct {
		Name     string `mapstructure:"name"`
		Version  string `mapstructure:"version"`
		LogLevel string `mapstructure:"logLevel"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	Services struct {
		Workers int `mapstructure:"task_workers"`
	}
)

func New() (*Config, error) {
	cfg := Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	for _, k := range v.AllKeys() {
		anyValue := v.Get(k)
		str, ok := anyValue.(string)
		if !ok {
			continue
		}

		replaced := os.ExpandEnv(str)
		v.Set(k, replaced)
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling file: %w", err))
	}

	return &cfg, nil
}
