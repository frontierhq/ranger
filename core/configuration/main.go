package configuration

import (
	"strings"

	"github.com/frontierdigital/ranger/core"
	"github.com/spf13/viper"
)

func LoadConfiguration() (configuration *core.Configuration, err error) {
	viper.SetEnvPrefix("RANGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.BindEnv("ado.pat")
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&configuration)

	return configuration, err
}
