package server

import (
	"challenge/pkg/config"
	pb "challenge/pkg/proto"
	"challenge/pkg/repository"
	"challenge/pkg/service"
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestChallengeService_MakeShortLink(t *testing.T) {
	type args struct {
		in *pb.Link
	}
	tests := []struct {
		name     string
		in       *pb.Link
		wantURL  string
		wantCode int
		wantErr  bool
	}{
		{
			name:     "OK",
			in:       &pb.Link{Data: "https://google.com"},
			wantURL:  "https://www.google.com/",
			wantCode: 301,
			wantErr:  false,
		},
		{
			name:     "invalid link",
			in:       &pb.Link{Data: "invalid link"},
			wantURL:  "",
			wantCode: 0,
			wantErr:  true,
		},
	}

	cfg, err := config.Load("../../configs/local")
	if err != nil {
		t.Errorf("Error load config: %s", err.Error())
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

	s := NewChallengeServer(service)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.MakeShortLink(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChallengeService.MakeShortLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			req, err := http.NewRequest("GET", got.GetData(), nil)
			if err != nil {
				t.Errorf("Error do request to: %s", got.GetData())
			}

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Error do request: %s", err)
			}

			url := res.Request.URL.String()
			code := res.Request.Response.StatusCode

			fmt.Println(code, url)

			if code != tt.wantCode {
				t.Errorf("ChallengeService.MakeShortLink() = %v, want %v", code, tt.wantCode)
			}

			if url != tt.wantURL {
				t.Errorf("ChallengeService.MakeShortLink() = %v, want %v", url, tt.wantURL)
			}
		})
	}
}
