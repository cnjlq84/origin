// Code generated by protoc-gen-go.
// source: health.proto
// DO NOT EDIT!

/*
Package grpc_health is a generated protocol buffer package.

It is generated from these files:
	health.proto

It has these top-level messages:
	HealthCheckRequest
	HealthCheckResponse
*/
package grpc_health

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type HealthCheckRequest struct {
}

func (m *HealthCheckRequest) Reset()         { *m = HealthCheckRequest{} }
func (m *HealthCheckRequest) String() string { return proto.CompactTextString(m) }
func (*HealthCheckRequest) ProtoMessage()    {}

type HealthCheckResponse struct {
}

func (m *HealthCheckResponse) Reset()         { *m = HealthCheckResponse{} }
func (m *HealthCheckResponse) String() string { return proto.CompactTextString(m) }
func (*HealthCheckResponse) ProtoMessage()    {}

func init() {
}

// Client API for HealthCheck service

type HealthCheckClient interface {
	Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type healthCheckClient struct {
	cc *grpc.ClientConn
}

func NewHealthCheckClient(cc *grpc.ClientConn) HealthCheckClient {
	return &healthCheckClient{cc}
}

func (c *healthCheckClient) Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := grpc.Invoke(ctx, "/grpc.health.HealthCheck/Check", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for HealthCheck service

type HealthCheckServer interface {
	Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

func RegisterHealthCheckServer(s *grpc.Server, srv HealthCheckServer) {
	s.RegisterService(&_HealthCheck_serviceDesc, srv)
}

func _HealthCheck_Check_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(HealthCheckServer).Check(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _HealthCheck_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.health.HealthCheck",
	HandlerType: (*HealthCheckServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _HealthCheck_Check_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
