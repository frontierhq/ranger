package config

import (
	"strings"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
	"github.com/spf13/viper"
)

func LoadConfig() (config *app.Config, err error) {
	viper.SetEnvPrefix("RANGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)

	return
}
