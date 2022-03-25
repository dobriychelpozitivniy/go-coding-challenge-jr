package server

import (
	pb "challenge/pkg/proto"
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func TestChallengeService_ReadMetadata(t *testing.T) {
	type args struct {
		md metadata.MD
		in *pb.Placeholder
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.Placeholder
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				md: metadata.Pairs("i-am-random-key", "random string"),
				in: &pb.Placeholder{},
			},
			want:    &pb.Placeholder{Data: "random string"},
			wantErr: false,
		},
		{
			name: "Empty i-am-random-key",
			args: args{
				md: nil,
				in: &pb.Placeholder{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	s := NewChallengeService(ChallengeServiceConfig{})

	go StartGRPCServer(s, ":8098")

	conn, err := grpc.Dial(":8098", grpc.WithInsecure())
	if err != nil {
		t.Errorf("Error dial grpc: %s", err)
	}

	client := pb.NewChallengeServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewOutgoingContext(context.Background(), tt.args.md)

			got, err := client.ReadMetadata(ctx, &pb.Placeholder{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ChallengeService.ReadMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got.ProtoReflect().Interface(), tt.want.ProtoReflect().Interface()) {
				t.Errorf("ChallengeService.ReadMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
