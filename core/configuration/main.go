package configuration

import (
	"strings"

	"github.com/frontierdigital/ranger/core/output"
	"github.com/spf13/viper"
)

type Configuration struct {
}

func LoadConfiguration() (configuration *Configuration, err error) {
	output.PrintlnLog("Loading configuration")

	viper.SetEnvPrefix("RANGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err = viper.Unmarshal(&configuration)

	return
}
