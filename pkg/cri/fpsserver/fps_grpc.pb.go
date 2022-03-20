// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package fpsserver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FPSServiceClient is the client API for FPSService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FPSServiceClient interface {
	FPSDrop(ctx context.Context, in *FPSDropRequest, opts ...grpc.CallOption) (*FPSDropReply, error)
}

type fPSServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFPSServiceClient(cc grpc.ClientConnInterface) FPSServiceClient {
	return &fPSServiceClient{cc}
}

func (c *fPSServiceClient) FPSDrop(ctx context.Context, in *FPSDropRequest, opts ...grpc.CallOption) (*FPSDropReply, error) {
	out := new(FPSDropReply)
	err := c.cc.Invoke(ctx, "/fpsserver.FPSService/FPSDrop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FPSServiceServer is the server API for FPSService service.
// All implementations must embed UnimplementedFPSServiceServer
// for forward compatibility
type FPSServiceServer interface {
	FPSDrop(context.Context, *FPSDropRequest) (*FPSDropReply, error)
	mustEmbedUnimplementedFPSServiceServer()
}

// UnimplementedFPSServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFPSServiceServer struct {
}

func (UnimplementedFPSServiceServer) FPSDrop(context.Context, *FPSDropRequest) (*FPSDropReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FPSDrop not implemented")
}
func (UnimplementedFPSServiceServer) mustEmbedUnimplementedFPSServiceServer() {}

// UnsafeFPSServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FPSServiceServer will
// result in compilation errors.
type UnsafeFPSServiceServer interface {
	mustEmbedUnimplementedFPSServiceServer()
}

func RegisterFPSServiceServer(s grpc.ServiceRegistrar, srv FPSServiceServer) {
	s.RegisterService(&FPSService_ServiceDesc, srv)
}

func _FPSService_FPSDrop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FPSDropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FPSServiceServer).FPSDrop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fpsserver.FPSService/FPSDrop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FPSServiceServer).FPSDrop(ctx, req.(*FPSDropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FPSService_ServiceDesc is the grpc.ServiceDesc for FPSService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FPSService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fpsserver.FPSService",
	HandlerType: (*FPSServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FPSDrop",
			Handler:    _FPSService_FPSDrop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fps.proto",
}