// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: multipaxos.proto

package multipaxos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MultiPaxosService_TriggerPrepare_FullMethodName = "/scooter.MultiPaxosService/TriggerPrepare"
	MultiPaxosService_Prepare_FullMethodName        = "/scooter.MultiPaxosService/Prepare"
	MultiPaxosService_Accept_FullMethodName         = "/scooter.MultiPaxosService/Accept"
	MultiPaxosService_Commit_FullMethodName         = "/scooter.MultiPaxosService/Commit"
)

// MultiPaxosServiceClient is the client API for MultiPaxosService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MultiPaxosServiceClient interface {
	TriggerPrepare(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error)
	Prepare(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error)
	Accept(ctx context.Context, in *AcceptRequest, opts ...grpc.CallOption) (*AcceptResponse, error)
	Commit(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*CommitResponse, error)
}

type multiPaxosServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMultiPaxosServiceClient(cc grpc.ClientConnInterface) MultiPaxosServiceClient {
	return &multiPaxosServiceClient{cc}
}

func (c *multiPaxosServiceClient) TriggerPrepare(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrepareResponse)
	err := c.cc.Invoke(ctx, MultiPaxosService_TriggerPrepare_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiPaxosServiceClient) Prepare(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrepareResponse)
	err := c.cc.Invoke(ctx, MultiPaxosService_Prepare_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiPaxosServiceClient) Accept(ctx context.Context, in *AcceptRequest, opts ...grpc.CallOption) (*AcceptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AcceptResponse)
	err := c.cc.Invoke(ctx, MultiPaxosService_Accept_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiPaxosServiceClient) Commit(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*CommitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommitResponse)
	err := c.cc.Invoke(ctx, MultiPaxosService_Commit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MultiPaxosServiceServer is the server API for MultiPaxosService service.
// All implementations must embed UnimplementedMultiPaxosServiceServer
// for forward compatibility.
type MultiPaxosServiceServer interface {
	TriggerPrepare(context.Context, *PrepareRequest) (*PrepareResponse, error)
	Prepare(context.Context, *PrepareRequest) (*PrepareResponse, error)
	Accept(context.Context, *AcceptRequest) (*AcceptResponse, error)
	Commit(context.Context, *CommitRequest) (*CommitResponse, error)
	mustEmbedUnimplementedMultiPaxosServiceServer()
}

// UnimplementedMultiPaxosServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMultiPaxosServiceServer struct{}

func (UnimplementedMultiPaxosServiceServer) TriggerPrepare(context.Context, *PrepareRequest) (*PrepareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerPrepare not implemented")
}
func (UnimplementedMultiPaxosServiceServer) Prepare(context.Context, *PrepareRequest) (*PrepareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepare not implemented")
}
func (UnimplementedMultiPaxosServiceServer) Accept(context.Context, *AcceptRequest) (*AcceptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Accept not implemented")
}
func (UnimplementedMultiPaxosServiceServer) Commit(context.Context, *CommitRequest) (*CommitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedMultiPaxosServiceServer) mustEmbedUnimplementedMultiPaxosServiceServer() {}
func (UnimplementedMultiPaxosServiceServer) testEmbeddedByValue()                           {}

// UnsafeMultiPaxosServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MultiPaxosServiceServer will
// result in compilation errors.
type UnsafeMultiPaxosServiceServer interface {
	mustEmbedUnimplementedMultiPaxosServiceServer()
}

func RegisterMultiPaxosServiceServer(s grpc.ServiceRegistrar, srv MultiPaxosServiceServer) {
	// If the following call pancis, it indicates UnimplementedMultiPaxosServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MultiPaxosService_ServiceDesc, srv)
}

func _MultiPaxosService_TriggerPrepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiPaxosServiceServer).TriggerPrepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MultiPaxosService_TriggerPrepare_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiPaxosServiceServer).TriggerPrepare(ctx, req.(*PrepareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiPaxosService_Prepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiPaxosServiceServer).Prepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MultiPaxosService_Prepare_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiPaxosServiceServer).Prepare(ctx, req.(*PrepareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiPaxosService_Accept_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiPaxosServiceServer).Accept(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MultiPaxosService_Accept_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiPaxosServiceServer).Accept(ctx, req.(*AcceptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiPaxosService_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiPaxosServiceServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MultiPaxosService_Commit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiPaxosServiceServer).Commit(ctx, req.(*CommitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MultiPaxosService_ServiceDesc is the grpc.ServiceDesc for MultiPaxosService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MultiPaxosService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "scooter.MultiPaxosService",
	HandlerType: (*MultiPaxosServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TriggerPrepare",
			Handler:    _MultiPaxosService_TriggerPrepare_Handler,
		},
		{
			MethodName: "Prepare",
			Handler:    _MultiPaxosService_Prepare_Handler,
		},
		{
			MethodName: "Accept",
			Handler:    _MultiPaxosService_Accept_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _MultiPaxosService_Commit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "multipaxos.proto",
}
