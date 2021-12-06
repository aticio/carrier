package binderserver

import (
	pb "carrier/Binder/rpc/binder"
	"errors"

	"golang.org/x/net/context"
)

type Server struct{}

func NewBinderServer() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) HandleServer(ctx context.Context, r *pb.Server) (*pb.Result, error) {
	seed, err := getServerSeed(r)
	if err != nil {
		return &pb.Result{
			Status: "Fail",
		}, nil
	}
	result := addSeed(seed)
	if !result {
		return &pb.Result{
			Status: "Fail",
		}, errors.New("Error adding seed")
	}
	return &pb.Result{
		Status: "Success",
	}, nil
}
