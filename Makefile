gen_proto:
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. pkg/proto/challenge.proto

run:
	go run cmd/server/main.go