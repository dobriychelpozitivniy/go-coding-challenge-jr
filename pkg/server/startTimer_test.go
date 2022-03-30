package server

import (
	"challenge/pkg/config"
	pb "challenge/pkg/proto"
	"challenge/pkg/repository"
	"challenge/pkg/service"
	"context"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestChallengeServer_StartTimer(t *testing.T) {
	cfg, err := config.Load("../../configs/local")
	if err != nil {
		t.Errorf("Error load config: %s", err.Error())
	}

	repo := repository.NewRepository(make(map[string][]repository.StreamChannels))

	service := service.NewService(service.ServiceConfig{
		ShortLinkServiceConfig: service.ShortLinkServiceConfig{
			BitlyURL:    cfg.BitlyURL,
			AccessToken: cfg.BitlyOAuthToken,
		},
		TimerServiceConfig: service.TimerServiceConfig{
			TimerURL: cfg.TimerURL,
		},
	}, repo)

	s := NewChallengeServer(service)

	go StartGRPCServer(s, ":8099")

	time.Sleep(time.Second * 3)

	testChallengeServer_StartTimer_OK(t)
	fmt.Println("First test Complete")
	testChallengeServer_StartTimer_Reconnect(t)
	testChallengeServer_StartTimer_TwoClients(t)
}

func testChallengeServer_StartTimer_OK(t *testing.T) {
	type args struct {
		timer  *pb.Timer
		stream pb.ChallengeService_StartTimerServer
	}
	tt := struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "OK",
		args:    args{timer: &pb.Timer{Seconds: 5, Frequency: 1, Name: "weaoijwd4oiawjdji"}},
		wantErr: false,
	}

	conn, err := grpc.Dial(":8099", grpc.WithInsecure())
	if err != nil {
		t.Errorf("Error dial grpc: %s", err)
	}

	client := pb.NewChallengeServiceClient(conn)

	t.Run(tt.name, func(t *testing.T) {
		timeout := time.After(time.Duration(tt.args.timer.Seconds+1) * time.Second)
		done := make(chan bool)
		go func() {
			select {
			case <-timeout:
				t.Fatal("Test didn't finish in time")
			case <-done:
				return
			}
		}()

		stream, err := client.StartTimer(context.Background(), tt.args.timer)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChallengeServer.StartTimer() error = %v, wantErr %v", err, tt.wantErr)
		}

		for {
			feature, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}
			log.Println("FEATURE: ", feature)
		}

		done <- true
	})
}

func testChallengeServer_StartTimer_Reconnect(t *testing.T) {
	type args struct {
		timer  *pb.Timer
		stream pb.ChallengeService_StartTimerServer
	}
	tt := struct {
		name            string
		args            args
		wantTicksBefore int
		wantTicksAfter  int
		sleepSeconds    int
		wantErr         bool
	}{
		name:            "RECONNECT",
		args:            args{timer: &pb.Timer{Seconds: 10, Frequency: 1, Name: "weaoijw1doiawjdji"}},
		wantTicksBefore: 4,
		wantTicksAfter:  3,
		sleepSeconds:    2,
		wantErr:         false,
	}

	conn, err := grpc.Dial(":8099", grpc.WithInsecure())
	if err != nil {
		t.Errorf("Error dial grpc: %s", err)
	}

	client := pb.NewChallengeServiceClient(conn)

	t.Run(tt.name, func(t *testing.T) {
		timeout := time.After(time.Duration(tt.args.timer.Seconds+1) * time.Second)
		done := make(chan bool)
		go func() {
			select {
			case <-timeout:
				t.Fatal("Test didn't finish in time")
			case <-done:
				return
			}
		}()

		stream, err := client.StartTimer(context.Background(), tt.args.timer)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChallengeServer.StartTimer() error = %v, wantErr %v", err, tt.wantErr)
		}

		var countTicks int = 0

		for {
			if tt.wantTicksBefore == countTicks {
				stream.CloseSend()
				break
			}

			feature, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}
			countTicks++
			log.Println("FEATURE: ", feature)
		}

		time.Sleep(time.Duration(tt.sleepSeconds) * time.Second)

		stream, err = client.StartTimer(context.Background(), tt.args.timer)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChallengeServer.StartTimer() error = %v, wantErr %v", err, tt.wantErr)
		}

		for {
			feature, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}

			log.Println("FEATURE: ", feature)
		}

		done <- true
	})
}

func testChallengeServer_StartTimer_TwoClients(t *testing.T) {
	type args struct {
		timer  *pb.Timer
		stream pb.ChallengeService_StartTimerServer
	}
	tt := struct {
		name              string
		args              args
		wantSecondClient  int
		startSecondClient int
		wantErr           bool
	}{
		name:              "two clients",
		args:              args{timer: &pb.Timer{Seconds: 10, Frequency: 1, Name: "wea2oijwdoiawjdji"}},
		wantSecondClient:  3,
		startSecondClient: 7,
		wantErr:           false,
	}

	conn, err := grpc.Dial(":8099", grpc.WithInsecure())
	if err != nil {
		t.Errorf("Error dial grpc: %s", err)
	}

	client := pb.NewChallengeServiceClient(conn)
	client2 := pb.NewChallengeServiceClient(conn)

	t.Run(tt.name, func(t *testing.T) {
		startSecondClient := time.After(time.Duration(tt.startSecondClient) * time.Second)
		go func() {
			<-startSecondClient
			go runSecondClient(t, client2, tt.wantErr, tt.args.timer, tt.wantSecondClient)
		}()

		stream, err := client.StartTimer(context.Background(), tt.args.timer)
		if (err != nil) != tt.wantErr {
			t.Errorf("ChallengeServer.StartTimer() error = %v, wantErr %v", err, tt.wantErr)
		}

		for {
			feature, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}

			log.Println("FIRST CLIENT TICK: ", feature.Seconds)
		}
	})
}

func runSecondClient(t *testing.T, client pb.ChallengeServiceClient, wantErr bool, timer *pb.Timer, wantSecondClient int) {
	timeout := time.After(time.Duration(wantSecondClient+1) * time.Second)
	done := make(chan bool)
	go func() {
		select {
		case <-timeout:
			t.Fatal("Test didn't finish in time")
		case <-done:
			return
		}
	}()

	stream, err := client.StartTimer(context.Background(), timer)
	if (err != nil) != wantErr {
		t.Errorf("ChallengeServer.StartTimer() error = %v, wantErr %v", err, wantErr)
	}

	var countTicks int = 0

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		countTicks++
		log.Println("SECOND CLIENT TICK: ", feature.Seconds)
	}

	done <- true
}
