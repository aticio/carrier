package carrierserver

import (
	pb "carrier/CarrierJvm/rpc/carrierjvm"

	b "carrier/CarrierJvm/domestic/business"

	"golang.org/x/net/context"
)

type Server struct {
}

//NewCarrierGrpcServer creates new seed grpc server
func NewCarrierJvmServer() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) UpdateCarrierJvm(ctx context.Context, r *pb.CarrierJvm) (*pb.Result, error) {
	err := b.UpdateSeed(r.GetHost(), r.GetPort())
	if err != nil {
		return &pb.Result{Status: "Fail"}, nil
	}
	return &pb.Result{Status: "Success"}, nil
}

func (s *Server) AddCarrierJvm(ctx context.Context, r *pb.CarrierJvm) (*pb.Result, error) {
	err := b.AddSeed(r.GetHost(), r.GetPort())
	if err != nil {
		return &pb.Result{Status: "Fail"}, nil
	}
	return &pb.Result{Status: "Success"}, nil
}

func (s *Server) DeleteCarrierJvm(ctx context.Context, r *pb.Void) (*pb.Result, error) {
	err := b.DeleteSeed()
	if err != nil {
		return &pb.Result{Status: "Fail"}, nil
	}
	return &pb.Result{Status: "Success"}, nil
}
