package carrierserver

import (
	pb "carrier/Carrier/rpc/carrier"

	b "carrier/Carrier/domestic/business"

	"golang.org/x/net/context"
)

type Server struct {
}

//NewCarrierGrpcServer creates new seed grpc server
func NewCarrierServer() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) UpdateCarrier(ctx context.Context, r *pb.Carrier) (*pb.Result, error) {
	err := b.UpdateSeed(r.GetHost(), r.GetPort())
	if err != nil {
		return &pb.Result{Status: "Fail"}, nil
	}
	return &pb.Result{Status: "Success"}, nil
}

func (s *Server) AddCarrier(ctx context.Context, r *pb.Carrier) (*pb.Result, error) {
	err := b.AddSeed(r.GetHost(), r.GetPort())
	if err != nil {
		return &pb.Result{Status: "Fail"}, nil
	}
	return &pb.Result{Status: "Success"}, nil
}

func (s *Server) DeleteCarrier(ctx context.Context, r *pb.Void) (*pb.Result, error) {
	err := b.DeleteSeed()
	if err != nil {
		return &pb.Result{Status: "Fail"}, nil
	}
	return &pb.Result{Status: "Success"}, nil
}
