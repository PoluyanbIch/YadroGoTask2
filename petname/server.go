package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	petname "github.com/dustinkirkland/golang-petname"

	petnamepb "yadro.com/course/proto"
)

type server struct {
	petnamepb.UnimplementedPetnameGeneratorServer
}

func (s *server) Ping(_ context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) Generate(ctx context.Context, r *petnamepb.PetnameRequest) (*petnamepb.PetnameResponse, error) {
	words := r.Words
	if words <= 0 {
		return nil, status.Error(codes.InvalidArgument, "words must be greater than 0")
	}
	separator := r.Separator
	name := petname.Generate(int(words), separator)
	return &petnamepb.PetnameResponse{Name: name}, nil
}

func (s *server) GenerateMany(req *petnamepb.PetnameStreamRequest, stream grpc.ServerStreamingServer[petnamepb.PetnameResponse]) error {
	words := req.Words
	separator := req.Separator
	namesCount := req.Names

	if words <= 0 {
		return status.Error(codes.InvalidArgument, "words must be greater than 0")
	}

	if namesCount <= 0 {
		return status.Error(codes.InvalidArgument, "names must be greater than 0")
	}

	for i := int64(0); i != namesCount; i++ {
		name := petname.Generate(int(words), separator)
		if err := stream.Send(&petnamepb.PetnameResponse{Name: name}); err != nil {
			return status.Error(codes.Unavailable, err.Error())
		}
	}
	return nil
}

func main() {
	var address string
	flag.StringVar(&address, "address", ":8080", "server address")
	flag.Parse()
	if envAddr := os.Getenv("PETNAME_GRPC_PORT"); envAddr != "" {
		address = ":" + envAddr
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	petnamepb.RegisterPetnameGeneratorServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
