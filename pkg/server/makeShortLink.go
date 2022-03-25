package server

import (
	"bytes"
	pb "challenge/pkg/proto"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type requestBody struct {
	LongUrl string `json:"long_url"`
}

type responseBody struct {
	Link string `json:"link"`
}

func (s *ChallengeService) MakeShortLink(ctx context.Context, in *pb.Link) (*pb.Link, error) {
	longUrl := in.GetData()

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	link, err := s.getShortLink(ctxTimeout, longUrl)
	if err != nil {
		return nil, err
	}

	res := &pb.Link{Data: link}

	return res, nil
}

func (s *ChallengeService) getShortLink(ctx context.Context, longUrl string) (string, error) {
	url := "https://api-ssl.bitly.com/v4/shorten"
	body := requestBody{LongUrl: longUrl}

	json_body, err := json.Marshal(body)
	if err != nil {
		return "", status.Error(codes.Internal, "Error marshal body")
	}

	bearer := "Bearer " + s.AccessToken

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_body))
	if err != nil {
		return "", status.Error(codes.Internal, "Error create request to bitly")
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", status.Error(codes.Internal, "Error do request to bitly")
	}

	var res responseBody

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", status.Error(codes.Internal, "Error decode response body from bitly")
	}

	return res.Link, nil
}
