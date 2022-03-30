package server

import (
	pb "challenge/pkg/proto"
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *ChallengeServer) ReadMetadata(ctx context.Context, in *pb.Placeholder) (*pb.Placeholder, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Empty metadata")
	}

	arrData := md.Get("i-am-random-key")

	data := strings.Join(arrData, ",")
	if data == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty i-am-random-key")
	}

	resp := &pb.Placeholder{Data: data}

	return resp, nil
}
