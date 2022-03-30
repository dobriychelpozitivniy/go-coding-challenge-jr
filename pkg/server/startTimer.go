package server

import (
	pb "challenge/pkg/proto"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChallengeServer) StartTimer(timer *pb.Timer, stream pb.ChallengeService_StartTimerServer) error {
	freq := timer.GetFrequency()
	seconds := timer.GetSeconds()
	name := timer.GetName()

	timerCh, doneCh, err := s.service.StartTimer(name, int(freq), int(seconds))
	if err != nil {
		fmt.Printf("Error start timer: %s", err)
		return status.Errorf(codes.Internal, "Error start timer: %s", err)
	}

	for {
		select {
		case <-doneCh:
			return nil
		case t := <-timerCh:
			stream.Send(t)
		}
	}
}
