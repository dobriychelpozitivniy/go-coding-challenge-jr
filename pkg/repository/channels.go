package repository

import (
	pb "challenge/pkg/proto"
	"sync"
)

type StreamChannels struct {
	Channel (chan *pb.Timer)
	Done    (chan bool)
}

type ChannelsRepository struct {
	channels map[string][]StreamChannels
	mu       sync.Mutex
}

func NewChannelsRepository(ch map[string][]StreamChannels) *ChannelsRepository {
	return &ChannelsRepository{channels: ch}
}

func (r *ChannelsRepository) GetChannel(name string) []StreamChannels {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.channels[name]
}

// func (r *ChannelsRepository) CheckEndChannel(name string) bool {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if streams, ok := r.channels[name]; ok {
// 		return streams.IsEnd
// 	}

// 	return true
// }

func (r *ChannelsRepository) AddChannel(name string) (chan *pb.Timer, chan bool, error) {
	r.mu.Lock()
	timerCh := make(chan *pb.Timer)
	doneCh := make(chan bool)
	streamChannels := StreamChannels{Channel: timerCh, Done: doneCh}

	streams, _ := r.channels[name]
	streams = append(streams, streamChannels)
	r.channels[name] = streams

	r.mu.Unlock()

	return timerCh, doneCh, nil
}

func (r *ChannelsRepository) CheckChannel(name string) bool {
	r.mu.Lock()
	_, ok := r.channels[name]
	r.mu.Unlock()

	return ok
}

func (r *ChannelsRepository) DeleteChannel(name string) {
	r.mu.Lock()
	delete(r.channels, name)
	r.mu.Unlock()
}

// func (r *ChannelsRepository) EndChannel(name string) {
// 	r.mu.Lock()
// 	if streams, ok := r.channels[name]; ok {
// 		streams.IsEnd = true
// 		r.channels[name] = streams
// 	}
// 	r.mu.Unlock()
// }
