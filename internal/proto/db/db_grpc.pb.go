// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: api/db.proto

package db

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	StorageEndpoint_Get_FullMethodName = "/jaeger_messages.StorageEndpoint/Get"
	StorageEndpoint_Set_FullMethodName = "/jaeger_messages.StorageEndpoint/Set"
	StorageEndpoint_Del_FullMethodName = "/jaeger_messages.StorageEndpoint/Del"
)

// StorageEndpointClient is the client API for StorageEndpoint service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageEndpointClient interface {
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*KeyValue, error)
	Set(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Del(ctx context.Context, in *Key, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type storageEndpointClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageEndpointClient(cc grpc.ClientConnInterface) StorageEndpointClient {
	return &storageEndpointClient{cc}
}

func (c *storageEndpointClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*KeyValue, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(KeyValue)
	err := c.cc.Invoke(ctx, StorageEndpoint_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageEndpointClient) Set(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, StorageEndpoint_Set_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageEndpointClient) Del(ctx context.Context, in *Key, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, StorageEndpoint_Del_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageEndpointServer is the server API for StorageEndpoint service.
// All implementations must embed UnimplementedStorageEndpointServer
// for forward compatibility.
type StorageEndpointServer interface {
	Get(context.Context, *Key) (*KeyValue, error)
	Set(context.Context, *KeyValue) (*emptypb.Empty, error)
	Del(context.Context, *Key) (*emptypb.Empty, error)
	mustEmbedUnimplementedStorageEndpointServer()
}

// UnimplementedStorageEndpointServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStorageEndpointServer struct{}

func (UnimplementedStorageEndpointServer) Get(context.Context, *Key) (*KeyValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedStorageEndpointServer) Set(context.Context, *KeyValue) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedStorageEndpointServer) Del(context.Context, *Key) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedStorageEndpointServer) mustEmbedUnimplementedStorageEndpointServer() {}
func (UnimplementedStorageEndpointServer) testEmbeddedByValue()                         {}

// UnsafeStorageEndpointServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageEndpointServer will
// result in compilation errors.
type UnsafeStorageEndpointServer interface {
	mustEmbedUnimplementedStorageEndpointServer()
}

func RegisterStorageEndpointServer(s grpc.ServiceRegistrar, srv StorageEndpointServer) {
	// If the following call pancis, it indicates UnimplementedStorageEndpointServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StorageEndpoint_ServiceDesc, srv)
}

func _StorageEndpoint_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageEndpointServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageEndpoint_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageEndpointServer).Get(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageEndpoint_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageEndpointServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageEndpoint_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageEndpointServer).Set(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageEndpoint_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageEndpointServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageEndpoint_Del_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageEndpointServer).Del(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

// StorageEndpoint_ServiceDesc is the grpc.ServiceDesc for StorageEndpoint service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StorageEndpoint_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jaeger_messages.StorageEndpoint",
	HandlerType: (*StorageEndpointServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _StorageEndpoint_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _StorageEndpoint_Set_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _StorageEndpoint_Del_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/db.proto",
}
