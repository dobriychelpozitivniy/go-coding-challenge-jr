package repository

import pb "challenge/pkg/proto"

type Channels interface {
	AddChannel(name string) (chan *pb.Timer, chan bool, error)
	CheckChannel(name string) bool
	DeleteChannel(name string)
	GetChannel(name string) []StreamChannels
}

type Repository struct {
	Channels
}

func NewRepository(ch map[string][]StreamChannels) *Repository {
	return &Repository{Channels: NewChannelsRepository(ch)}
}
