package config

import (
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	BitlyOAuthToken string `mapstructure:"BITLY_OAUTH_TOKEN"`
}

func Load(cfgPath string) (*Config, error) {
	var config Config

	viper.AddConfigPath(path.Dir(cfgPath))
	viper.SetConfigName(path.Base(cfgPath))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
