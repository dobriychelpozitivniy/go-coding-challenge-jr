package server

import (
	pb "challenge/pkg/proto"
	"fmt"
	"net"

	middlewareLogger "challenge/pkg/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ChallengeServiceConfig struct {
	AccessToken string
}

type ChallengeService struct {
	pb.UnimplementedChallengeServiceServer
	ChallengeServiceConfig
}

// func (s *Server) StartTimer(*Timer, ChallengeService_StartTimerServer) error {
// 	return status.Errorf(codes.Unimplemented, "method StartTimer not implemented")
// }

func NewChallengeService(cfg ChallengeServiceConfig) *ChallengeService {
	return &ChallengeService{ChallengeServiceConfig: cfg}
}

// Init grpc server and start him
func StartGRPCServer(server *ChallengeService, port string) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic().Msgf("Error in listen %s", err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middlewareLogger.NewUnaryServerInterceptor(),
		)),
	)

	pb.RegisterChallengeServiceServer(grpcServer, server)

	fmt.Println("Start server")

	if err := grpcServer.Serve(listen); err != nil {
		log.Panic().Msgf("error serve grpc: %s", err.Error())
	}
}
