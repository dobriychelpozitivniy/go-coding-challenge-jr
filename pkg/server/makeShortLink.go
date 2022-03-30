package server

import (
	pb "challenge/pkg/proto"
	"context"
	"time"
)

func (s *ChallengeServer) MakeShortLink(ctx context.Context, in *pb.Link) (*pb.Link, error) {
	longUrl := in.GetData()

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	link, err := s.service.GetShortLink(ctxTimeout, longUrl)
	if err != nil {
		return nil, err
	}

	res := &pb.Link{Data: link}

	return res, nil
}
