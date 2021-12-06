package seedserver

import (
	"carrier/Seed/domestic/dao"
	pb "carrier/Seed/rpc/seed"
	"fmt"

	"golang.org/x/net/context"

	"github.com/spf13/viper"
)

type Server struct {
	db dao.SeedDAO
}

//NewSeedGrpcServer creates new seed grpc server
func NewSeedServer() (*Server, error) {
	handler, err := dao.NewSeedDAO(viper.GetString("mongoUrl"))
	if err != nil {
		return nil, fmt.Errorf("could not ceate a database handler object, error %v", err)
	}
	return &Server{
		db: *handler,
	}, nil
}

func (s *Server) AddSeed(ctx context.Context, r *pb.Seed) (*pb.Result, error) {
	seed := convertPbToSeed(r)
	err := s.db.AddSeed(*seed)

	if err != nil {
		return &pb.Result{
			Status: "Fail",
		}, nil
	}
	return &pb.Result{
		Status: "Successful",
	}, nil
}

func (s *Server) GetSeed(ctx context.Context, r *pb.Seed) (*pb.Seed, error) {
	seed, err := s.db.GetSeed(r.GetHost(), r.GetPort())
	if err != nil {
		return nil, err
	}
	pbseed := convertSeedToPb(seed)
	return pbseed, nil
}

func (s *Server) GetAll(ctx context.Context, r *pb.Void) (*pb.Seeds, error) {
	seeds, err := s.db.GetAll()
	if err != nil {
		return nil, err
	}
	return convertSeedsToPb(seeds), nil
}

func (s *Server) UpdateSeed(ctx context.Context, r *pb.Seed) (*pb.Result, error) {
	err := s.db.UpdateSeed(*convertPbToSeed(r))
	if err != nil {
		return &pb.Result{
			Status: "Fail",
		}, err
	}
	return &pb.Result{
		Status: "Successful",
	}, nil
}
