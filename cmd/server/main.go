package main

import (
	"challenge/pkg/config"
	"challenge/pkg/server"
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	flag.String("config", "configs/local", "display colorized output")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Panic().Msgf("Error parse flag config path: %s", err)
	}

	cfg, err := config.Load(viper.GetString("config"))
	if err != nil {
		log.Panic().Msgf("Error init config: %s", err.Error())
	}

	s := server.NewChallengeService(server.ChallengeServiceConfig{
		AccessToken: cfg.BitlyOAuthToken,
	})

	server.StartGRPCServer(s, cfg.Port)
}
