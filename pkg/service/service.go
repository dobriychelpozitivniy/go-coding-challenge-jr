package service

import (
	pb "challenge/pkg/proto"
	"challenge/pkg/repository"
	"context"
)

type ShortLink interface {
	GetShortLink(ctx context.Context, longUrl string) (string, error)
}

type Timer interface {
	StartTimer(name string, frequency int, seconds int) (chan *pb.Timer, chan bool, error)
	CheckExistTimer(name string) bool
	CloseTimer(name string)
	SendToChannels(name string, frequency int64, seconds int64)
	StartTicker(name string, frequency int64, seconds int64)
}

type Service struct {
	ShortLink
	Timer
}

type ServiceConfig struct {
	ShortLinkServiceConfig
	TimerServiceConfig
}

func NewService(cfg ServiceConfig, r *repository.Repository) *Service {
	return &Service{
		Timer:     NewTimerService(cfg.TimerServiceConfig, r.Channels),
		ShortLink: NewShortLinkService(cfg.ShortLinkServiceConfig),
	}
}
