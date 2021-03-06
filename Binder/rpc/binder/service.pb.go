// Code generated by protoc-gen-go. DO NOT EDIT.
// seed: service.proto

/*
Package com_carrier_server is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	Server
	Result
*/
package com_carrier_server

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Server struct {
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Port string `protobuf:"bytes,2,opt,name=port" json:"port,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Server) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Server) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

type Result struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*Server)(nil), "com.carrier.server.Server")
	proto.RegisterType((*Result)(nil), "com.carrier.server.Result")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ServerService service

type ServerServiceClient interface {
	HandleServer(ctx context.Context, in *Server, opts ...grpc.CallOption) (*Result, error)
}

type serverServiceClient struct {
	cc *grpc.ClientConn
}

func NewServerServiceClient(cc *grpc.ClientConn) ServerServiceClient {
	return &serverServiceClient{cc}
}

func (c *serverServiceClient) HandleServer(ctx context.Context, in *Server, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/com.carrier.server.ServerService/HandleServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ServerService service

type ServerServiceServer interface {
	HandleServer(context.Context, *Server) (*Result, error)
}

func RegisterServerServiceServer(s *grpc.Server, srv ServerServiceServer) {
	s.RegisterService(&_ServerService_serviceDesc, srv)
}

func _ServerService_HandleServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Server)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).HandleServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.carrier.server.ServerService/HandleServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).HandleServer(ctx, req.(*Server))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.carrier.server.ServerService",
	HandlerType: (*ServerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleServer",
			Handler:    _ServerService_HandleServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x4a, 0xce, 0xcf, 0xd5, 0x4b,
	0xce, 0x2f, 0x2d, 0xca, 0x4c, 0x2d, 0xd2, 0x03, 0x49, 0xa5, 0x16, 0x29, 0x19, 0x70, 0xb1, 0x05,
	0x83, 0x59, 0x42, 0x42, 0x5c, 0x2c, 0x19, 0xf9, 0xc5, 0x25, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c,
	0x41, 0x60, 0x36, 0x48, 0xac, 0x20, 0xbf, 0xa8, 0x44, 0x82, 0x09, 0x22, 0x06, 0x62, 0x2b, 0x29,
	0x70, 0xb1, 0x05, 0xa5, 0x16, 0x97, 0xe6, 0x94, 0x08, 0x89, 0x71, 0xb1, 0x15, 0x97, 0x24, 0x96,
	0x94, 0x16, 0x43, 0xf5, 0x40, 0x79, 0x46, 0xe1, 0x5c, 0xbc, 0x10, 0x33, 0x83, 0x21, 0xd6, 0x0b,
	0xb9, 0x71, 0xf1, 0x78, 0x24, 0xe6, 0xa5, 0xe4, 0xa4, 0x42, 0xad, 0x92, 0xd2, 0xc3, 0x74, 0x89,
	0x1e, 0x44, 0x4e, 0x0a, 0xab, 0x1c, 0xc4, 0xc2, 0x24, 0x36, 0xb0, 0x3f, 0x8c, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xf5, 0x05, 0x44, 0x88, 0xd8, 0x00, 0x00, 0x00,
}
