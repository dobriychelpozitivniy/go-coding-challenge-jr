package config

import (
	"flag"
	"fmt"
	"path"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	BitlyOAuthToken string `mapstructure:"BITLY_OAUTH_TOKEN"`
	BitlyURL        string `mapstructure:"BITLY_URL"`
	TimerURL        string `mapstructure:"TIMER_URL"`
}

func LoadFlags() error {
	flag.String("config", "configs/local", "display colorized output")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return fmt.Errorf("Error parse flags: %s", err)
	}

	return nil
}

func Load(cfgPath string) (*Config, error) {
	var config Config

	viper.AddConfigPath(path.Dir(cfgPath))
	viper.SetConfigName(path.Base(cfgPath))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	viper.BindEnv("BITLY_OAUTH_TOKEN")

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
