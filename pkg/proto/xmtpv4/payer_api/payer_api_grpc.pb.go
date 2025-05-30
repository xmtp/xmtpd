// Payer API

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: xmtpv4/payer_api/payer_api.proto

package payer_api

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

const (
	PayerApi_PublishClientEnvelopes_FullMethodName = "/xmtp.xmtpv4.payer_api.PayerApi/PublishClientEnvelopes"
	PayerApi_GetReaderNode_FullMethodName          = "/xmtp.xmtpv4.payer_api.PayerApi/GetReaderNode"
)

// PayerApiClient is the client API for PayerApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PayerApiClient interface {
	// Publish envelope
	PublishClientEnvelopes(ctx context.Context, in *PublishClientEnvelopesRequest, opts ...grpc.CallOption) (*PublishClientEnvelopesResponse, error)
	GetReaderNode(ctx context.Context, in *GetReaderNodeRequest, opts ...grpc.CallOption) (*GetReaderNodeResponse, error)
}

type payerApiClient struct {
	cc grpc.ClientConnInterface
}

func NewPayerApiClient(cc grpc.ClientConnInterface) PayerApiClient {
	return &payerApiClient{cc}
}

func (c *payerApiClient) PublishClientEnvelopes(ctx context.Context, in *PublishClientEnvelopesRequest, opts ...grpc.CallOption) (*PublishClientEnvelopesResponse, error) {
	out := new(PublishClientEnvelopesResponse)
	err := c.cc.Invoke(ctx, PayerApi_PublishClientEnvelopes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *payerApiClient) GetReaderNode(ctx context.Context, in *GetReaderNodeRequest, opts ...grpc.CallOption) (*GetReaderNodeResponse, error) {
	out := new(GetReaderNodeResponse)
	err := c.cc.Invoke(ctx, PayerApi_GetReaderNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PayerApiServer is the server API for PayerApi service.
// All implementations must embed UnimplementedPayerApiServer
// for forward compatibility
type PayerApiServer interface {
	// Publish envelope
	PublishClientEnvelopes(context.Context, *PublishClientEnvelopesRequest) (*PublishClientEnvelopesResponse, error)
	GetReaderNode(context.Context, *GetReaderNodeRequest) (*GetReaderNodeResponse, error)
	mustEmbedUnimplementedPayerApiServer()
}

// UnimplementedPayerApiServer must be embedded to have forward compatible implementations.
type UnimplementedPayerApiServer struct {
}

func (UnimplementedPayerApiServer) PublishClientEnvelopes(context.Context, *PublishClientEnvelopesRequest) (*PublishClientEnvelopesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishClientEnvelopes not implemented")
}
func (UnimplementedPayerApiServer) GetReaderNode(context.Context, *GetReaderNodeRequest) (*GetReaderNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReaderNode not implemented")
}
func (UnimplementedPayerApiServer) mustEmbedUnimplementedPayerApiServer() {}

// UnsafePayerApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PayerApiServer will
// result in compilation errors.
type UnsafePayerApiServer interface {
	mustEmbedUnimplementedPayerApiServer()
}

func RegisterPayerApiServer(s grpc.ServiceRegistrar, srv PayerApiServer) {
	s.RegisterService(&PayerApi_ServiceDesc, srv)
}

func _PayerApi_PublishClientEnvelopes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishClientEnvelopesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PayerApiServer).PublishClientEnvelopes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PayerApi_PublishClientEnvelopes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PayerApiServer).PublishClientEnvelopes(ctx, req.(*PublishClientEnvelopesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PayerApi_GetReaderNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReaderNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PayerApiServer).GetReaderNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PayerApi_GetReaderNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PayerApiServer).GetReaderNode(ctx, req.(*GetReaderNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PayerApi_ServiceDesc is the grpc.ServiceDesc for PayerApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PayerApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xmtp.xmtpv4.payer_api.PayerApi",
	HandlerType: (*PayerApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishClientEnvelopes",
			Handler:    _PayerApi_PublishClientEnvelopes_Handler,
		},
		{
			MethodName: "GetReaderNode",
			Handler:    _PayerApi_GetReaderNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "xmtpv4/payer_api/payer_api.proto",
}
