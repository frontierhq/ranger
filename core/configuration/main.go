package configuration

import (
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
}

func LoadConfiguration() (configuration *Configuration, err error) {
	viper.SetEnvPrefix("RANGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err = viper.Unmarshal(&configuration)

	return
}
