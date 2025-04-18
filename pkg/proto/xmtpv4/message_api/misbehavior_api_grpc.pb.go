// API for reporting and querying node misbehavior in decentralized XMTP

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: xmtpv4/message_api/misbehavior_api.proto

package message_api

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
	MisbehaviorApi_SubmitMisbehaviorReport_FullMethodName = "/xmtp.xmtpv4.message_api.MisbehaviorApi/SubmitMisbehaviorReport"
	MisbehaviorApi_QueryMisbehaviorReports_FullMethodName = "/xmtp.xmtpv4.message_api.MisbehaviorApi/QueryMisbehaviorReports"
)

// MisbehaviorApiClient is the client API for MisbehaviorApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MisbehaviorApiClient interface {
	SubmitMisbehaviorReport(ctx context.Context, in *SubmitMisbehaviorReportRequest, opts ...grpc.CallOption) (*SubmitMisbehaviorReportResponse, error)
	QueryMisbehaviorReports(ctx context.Context, in *QueryMisbehaviorReportsRequest, opts ...grpc.CallOption) (*QueryMisbehaviorReportsResponse, error)
}

type misbehaviorApiClient struct {
	cc grpc.ClientConnInterface
}

func NewMisbehaviorApiClient(cc grpc.ClientConnInterface) MisbehaviorApiClient {
	return &misbehaviorApiClient{cc}
}

func (c *misbehaviorApiClient) SubmitMisbehaviorReport(ctx context.Context, in *SubmitMisbehaviorReportRequest, opts ...grpc.CallOption) (*SubmitMisbehaviorReportResponse, error) {
	out := new(SubmitMisbehaviorReportResponse)
	err := c.cc.Invoke(ctx, MisbehaviorApi_SubmitMisbehaviorReport_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *misbehaviorApiClient) QueryMisbehaviorReports(ctx context.Context, in *QueryMisbehaviorReportsRequest, opts ...grpc.CallOption) (*QueryMisbehaviorReportsResponse, error) {
	out := new(QueryMisbehaviorReportsResponse)
	err := c.cc.Invoke(ctx, MisbehaviorApi_QueryMisbehaviorReports_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MisbehaviorApiServer is the server API for MisbehaviorApi service.
// All implementations must embed UnimplementedMisbehaviorApiServer
// for forward compatibility
type MisbehaviorApiServer interface {
	SubmitMisbehaviorReport(context.Context, *SubmitMisbehaviorReportRequest) (*SubmitMisbehaviorReportResponse, error)
	QueryMisbehaviorReports(context.Context, *QueryMisbehaviorReportsRequest) (*QueryMisbehaviorReportsResponse, error)
	mustEmbedUnimplementedMisbehaviorApiServer()
}

// UnimplementedMisbehaviorApiServer must be embedded to have forward compatible implementations.
type UnimplementedMisbehaviorApiServer struct {
}

func (UnimplementedMisbehaviorApiServer) SubmitMisbehaviorReport(context.Context, *SubmitMisbehaviorReportRequest) (*SubmitMisbehaviorReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitMisbehaviorReport not implemented")
}
func (UnimplementedMisbehaviorApiServer) QueryMisbehaviorReports(context.Context, *QueryMisbehaviorReportsRequest) (*QueryMisbehaviorReportsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryMisbehaviorReports not implemented")
}
func (UnimplementedMisbehaviorApiServer) mustEmbedUnimplementedMisbehaviorApiServer() {}

// UnsafeMisbehaviorApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MisbehaviorApiServer will
// result in compilation errors.
type UnsafeMisbehaviorApiServer interface {
	mustEmbedUnimplementedMisbehaviorApiServer()
}

func RegisterMisbehaviorApiServer(s grpc.ServiceRegistrar, srv MisbehaviorApiServer) {
	s.RegisterService(&MisbehaviorApi_ServiceDesc, srv)
}

func _MisbehaviorApi_SubmitMisbehaviorReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitMisbehaviorReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MisbehaviorApiServer).SubmitMisbehaviorReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MisbehaviorApi_SubmitMisbehaviorReport_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MisbehaviorApiServer).SubmitMisbehaviorReport(ctx, req.(*SubmitMisbehaviorReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MisbehaviorApi_QueryMisbehaviorReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMisbehaviorReportsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MisbehaviorApiServer).QueryMisbehaviorReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MisbehaviorApi_QueryMisbehaviorReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MisbehaviorApiServer).QueryMisbehaviorReports(ctx, req.(*QueryMisbehaviorReportsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MisbehaviorApi_ServiceDesc is the grpc.ServiceDesc for MisbehaviorApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MisbehaviorApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xmtp.xmtpv4.message_api.MisbehaviorApi",
	HandlerType: (*MisbehaviorApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitMisbehaviorReport",
			Handler:    _MisbehaviorApi_SubmitMisbehaviorReport_Handler,
		},
		{
			MethodName: "QueryMisbehaviorReports",
			Handler:    _MisbehaviorApi_QueryMisbehaviorReports_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "xmtpv4/message_api/misbehavior_api.proto",
}
