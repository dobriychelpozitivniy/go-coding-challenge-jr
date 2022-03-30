package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type linkRequestBody struct {
	LongUrl string `json:"long_url"`
}

type linkResponseBody struct {
	Link string `json:"link"`
}

type ShortLinkServiceConfig struct {
	BitlyURL    string
	AccessToken string
}

type ShortLinkService struct {
	cfg ShortLinkServiceConfig
}

func NewShortLinkService(config ShortLinkServiceConfig) *ShortLinkService {
	return &ShortLinkService{cfg: config}
}

func (s *ShortLinkService) GetShortLink(ctx context.Context, longUrl string) (string, error) {
	url := s.cfg.BitlyURL
	body := linkRequestBody{LongUrl: longUrl}

	json_body, err := json.Marshal(body)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Error marshal body: %s", err)
	}

	bearer := "Bearer " + s.cfg.AccessToken

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_body))
	if err != nil {
		return "", status.Errorf(codes.Internal, "Error create request to bitly: %s", err)
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode > 400 {
		return "", status.Errorf(codes.Internal, "Error do request to bitly: %s", err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	if resp.StatusCode == 400 {
		return "", status.Errorf(codes.InvalidArgument, "Error do request to bitly: %s", err)
	}

	var res linkResponseBody

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Error decode response body from bitly: %s", err)
	}

	return res.Link, nil
}
