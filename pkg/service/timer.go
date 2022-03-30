package service

import (
	pb "challenge/pkg/proto"
	"challenge/pkg/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TimerServiceConfig struct {
	TimerURL string
}

type TimerService struct {
	cfg  TimerServiceConfig
	repo repository.Channels
}

type createTimerResponse struct {
	Timer string `json:"timer"`
}

type getTimerResponse struct {
	Timer            string  `json:"timer"`
	SecondsRemaining float64 `json:"seconds_remaining"`
}

func NewTimerService(config TimerServiceConfig, r repository.Channels) *TimerService {
	return &TimerService{cfg: config, repo: r}
}

// Check is timer exist
// If not exist - create timer with channels in local repository, create timer in timercheck.io API and run ticker
// If exist - create channels in local repository
func (s *TimerService) StartTimer(name string, frequency int, seconds int) (chan *pb.Timer, chan bool, error) {
	start := time.Now()
	if s.repo.CheckChannel(name) {
		timerCh, doneCh, err := s.repo.AddChannel(name)
		if err != nil {
			return nil, nil, err
		}

		return timerCh, doneCh, nil
	}
	fmt.Println("TIME:", time.Since(start))

	timerCh, doneCh, err := s.repo.AddChannel(name)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.CreateTimer(name, seconds)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Error create timer: %s", err.Error())
	}

	go s.StartTicker(name, int64(frequency), int64(seconds))

	return timerCh, doneCh, nil
}

// Create ticker with frequency and seconds, get timer from timercheck.io API, get stream's channels on timer's name from local repository and send in him timer info
// Where ticker is done, he send done in stream's done channels for close stream and delete timer from local repository
func (s *TimerService) StartTicker(name string, frequency int64, seconds int64) {
	ticker := time.NewTicker(time.Second * time.Duration(frequency))
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(time.Second * time.Duration(seconds))
		done <- true
	}()

	for {
		select {
		case <-done:
			go s.CloseTimer(name)
			return
		case <-ticker.C:
			res, err := s.GetTimer(name)

			if res == nil && err == nil {
				go s.CloseTimer(name)
				return
			}

			fmt.Println("RES", res)
			if err != nil {
				fmt.Printf("err get timer: %s", err)
				go s.CloseTimer(name)
				return
			}

			go s.SendToChannels(name, frequency, int64(res.SecondsRemaining))
		}
	}
}

func (s *TimerService) SendToChannels(name string, frequency int64, seconds int64) {
	channels := s.repo.GetChannel(name)

	for _, channel := range channels {
		channel.Channel <- &pb.Timer{Seconds: seconds, Frequency: frequency, Name: name}
	}
}

// Send true in done channels and delete timer from local repository
func (s *TimerService) CloseTimer(name string) {
	channels := s.repo.GetChannel(name)

	for _, channel := range channels {
		channel.Done <- true
	}

	s.repo.DeleteChannel(name)
}

// Check is exist timer with this name in local repository
func (s *TimerService) CheckExistTimer(name string) bool {
	return s.repo.CheckChannel(name)
}

// Create timer in https://timercheck.io API
func (s *TimerService) CreateTimer(name string, seconds int) (*createTimerResponse, error) {
	url := fmt.Sprintf("%s/%s/%v", s.cfg.TimerURL, name, seconds)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error create request to bitly: %s", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, status.Errorf(codes.Internal, "Error do request to %s: %s", url, err)
	}

	var res createTimerResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error decode response body from bitly: %s", err)
	}

	fmt.Println("CREATE", res)

	return &res, nil
}

// Get timer info from https://timercheck.io API
func (s *TimerService) GetTimer(name string) (*getTimerResponse, error) {
	url := fmt.Sprintf("%s/%s", s.cfg.TimerURL, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error create request to bitly: %s", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp.StatusCode == 504 {
		return nil, nil
	}

	if err != nil || resp.StatusCode != 200 {
		return nil, status.Errorf(codes.Internal, "Error do request to %s: %s", url, err)
	}

	var res getTimerResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error decode response body from bitly: %s", err)
	}

	fmt.Println("GET", res)

	return &res, nil
}
