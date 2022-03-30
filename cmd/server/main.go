package main

import (
	"challenge/pkg/config"
	"challenge/pkg/repository"
	"challenge/pkg/server"
	"challenge/pkg/service"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	err := config.LoadFlags()
	if err != nil {
		log.Panic().Msgf("Error load flags: %s", err.Error())
	}

	cfg, err := config.Load(viper.GetString("config"))
	if err != nil {
		log.Panic().Msgf("Error init config: %s", err.Error())
	}

	repo := repository.NewRepository(make(map[string][]repository.StreamChannels))

	service := service.NewService(service.ServiceConfig{
		ShortLinkServiceConfig: service.ShortLinkServiceConfig{
			BitlyURL:    *cfg.BitlyURL,
			AccessToken: *cfg.BitlyOAuthToken,
		},
		TimerServiceConfig: service.TimerServiceConfig{
			TimerURL: *cfg.TimerURL,
		},
	}, repo)

	s := server.NewChallengeServer(service)

	server.StartGRPCServer(s, *cfg.Port)
}
