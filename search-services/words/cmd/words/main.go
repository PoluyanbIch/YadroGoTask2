package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"yadro.com/course/words/internal/config"
	"yadro.com/course/words/internal/transport"

	pb "yadro.com/course/proto/words"
)

func main() {
	cfg := config.Load()
	listener, err := net.Listen("tcp", cfg.GRPCPORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWordsServer(s, transport.NewServer())
	reflection.Register(s)

	log.Printf("Words service listening on %s", cfg.GRPCPORT)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
