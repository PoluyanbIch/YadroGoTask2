package transport

import (
	"context"
	"log"

	"yadro.com/course/words/internal/app"
	"yadro.com/course/words/internal/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "yadro.com/course/proto/words"
)

type server struct {
	pb.UnimplementedWordsServer
}

func NewServer() pb.WordsServer {
	return &server{}
}

func (s *server) Norm(ctx context.Context, r *pb.WordsRequest) (*pb.WordsReply, error) {
	str := r.Phrase
	if len(str) > 4*1024 {
		return nil, status.Error(codes.ResourceExhausted, "message > 4 KiB")
	}
	banWords, err := config.LoadStopWords("/config/stopwords.txt")
	if err != nil {
		log.Fatalf("Failed to load stopwords %v", err)
	}
	n := app.NewNormalizer(banWords)

	norm := n.Normalize(str)

	return &pb.WordsReply{Words: norm}, nil
}

func (s *server) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
