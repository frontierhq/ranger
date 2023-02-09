package config

type ADOConfig struct {
	PAT string `mapstructure:"PAT"`
}

type GitConfig struct {
	UserEmail string `mapstructure:"UserEmail"`
	UserName  string `mapstructure:"UserName"`
}

type Config struct {
	ADO ADOConfig `mapstructure:"ADO"`
	Git GitConfig `mapstructure:"Git"`
}
