package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ADO ADOConfig `mapstructure:"ADO"`
	Git GitConfig `mapstructure:"Git"`
}

type ADOConfig struct {
	PAT string `mapstructure:"PAT"`
}

type GitConfig struct {
	UserEmail string `mapstructure:"UserEmail"`
	UserName  string `mapstructure:"UserName"`
}

func LoadConfig() (config *Config, err error) {
	viper.SetEnvPrefix("RANGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.BindEnv("ado.pat")
	if err != nil {
		return nil, err
	}

	err = viper.BindEnv("git.useremail")
	if err != nil {
		return nil, err
	}

	err = viper.BindEnv("git.username")
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)

	return config, err
}
