package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	DatabaseHost        string `mapstructure:"DATABASE_HOST"`
	DatabaseName        string `mapstructure:"DATABASE_NAME"`
	DatabaseUser        string `mapstructure:"DATABASE_USER"`
	DatabasePassword    string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseOpenTimeout int    `mapstructure:"DATABASE_OPEN_TIMEOUT"`
}

func InitializeConfig(configPath string) Configuration {
	// Location of the config
	viper.AddConfigPath(configPath)
	// Name of the config file
	viper.SetConfigName("config")
	// Extension of the config file
	viper.SetConfigType("env")

	// Override vales defined in the config with environmental variables (if present)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to read config %s\n", err))
	}

	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to initialize config %s\n", err))
	}

	return config
}
