package server

import (
	pb "challenge/pkg/proto"
	"fmt"
	"net"

	middlewareLogger "challenge/pkg/logger"

	"challenge/pkg/service"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ChallengeServer struct {
	pb.UnimplementedChallengeServiceServer
	service *service.Service
}

// func (s *Server) StartTimer(*Timer, ChallengeServer_StartTimerServer) error {
// 	return status.Errorf(codes.Unimplemented, "method StartTimer not implemented")
// }

func NewChallengeServer(s *service.Service) *ChallengeServer {
	return &ChallengeServer{service: s}
}

// Init grpc server and start him
func StartGRPCServer(server *ChallengeServer, port string) {
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
