// Message API

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mls_validation/v1/service.proto

package mls_validationv1

import (
	context "context"
	v1 "github.com/xmtp/xmtpd/pkg/proto/identity/api/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ValidationApi_ValidateGroupMessages_FullMethodName               = "/xmtp.mls_validation.v1.ValidationApi/ValidateGroupMessages"
	ValidationApi_GetAssociationState_FullMethodName                 = "/xmtp.mls_validation.v1.ValidationApi/GetAssociationState"
	ValidationApi_ValidateInboxIdKeyPackages_FullMethodName          = "/xmtp.mls_validation.v1.ValidationApi/ValidateInboxIdKeyPackages"
	ValidationApi_VerifySmartContractWalletSignatures_FullMethodName = "/xmtp.mls_validation.v1.ValidationApi/VerifySmartContractWalletSignatures"
)

// ValidationApiClient is the client API for ValidationApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ValidationApiClient interface {
	// Validates and parses a group message and returns relevant details
	ValidateGroupMessages(ctx context.Context, in *ValidateGroupMessagesRequest, opts ...grpc.CallOption) (*ValidateGroupMessagesResponse, error)
	// Gets the final association state for a batch of identity updates
	GetAssociationState(ctx context.Context, in *GetAssociationStateRequest, opts ...grpc.CallOption) (*GetAssociationStateResponse, error)
	// Validates InboxID key packages and returns credential information for them, without checking
	// whether an InboxId <> InstallationPublicKey pair is really valid.
	ValidateInboxIdKeyPackages(ctx context.Context, in *ValidateKeyPackagesRequest, opts ...grpc.CallOption) (*ValidateInboxIdKeyPackagesResponse, error)
	// Verifies smart contracts
	// This request is proxied from the node, so we'll reuse those messages.
	VerifySmartContractWalletSignatures(ctx context.Context, in *v1.VerifySmartContractWalletSignaturesRequest, opts ...grpc.CallOption) (*v1.VerifySmartContractWalletSignaturesResponse, error)
}

type validationApiClient struct {
	cc grpc.ClientConnInterface
}

func NewValidationApiClient(cc grpc.ClientConnInterface) ValidationApiClient {
	return &validationApiClient{cc}
}

func (c *validationApiClient) ValidateGroupMessages(ctx context.Context, in *ValidateGroupMessagesRequest, opts ...grpc.CallOption) (*ValidateGroupMessagesResponse, error) {
	out := new(ValidateGroupMessagesResponse)
	err := c.cc.Invoke(ctx, ValidationApi_ValidateGroupMessages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *validationApiClient) GetAssociationState(ctx context.Context, in *GetAssociationStateRequest, opts ...grpc.CallOption) (*GetAssociationStateResponse, error) {
	out := new(GetAssociationStateResponse)
	err := c.cc.Invoke(ctx, ValidationApi_GetAssociationState_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *validationApiClient) ValidateInboxIdKeyPackages(ctx context.Context, in *ValidateKeyPackagesRequest, opts ...grpc.CallOption) (*ValidateInboxIdKeyPackagesResponse, error) {
	out := new(ValidateInboxIdKeyPackagesResponse)
	err := c.cc.Invoke(ctx, ValidationApi_ValidateInboxIdKeyPackages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *validationApiClient) VerifySmartContractWalletSignatures(ctx context.Context, in *v1.VerifySmartContractWalletSignaturesRequest, opts ...grpc.CallOption) (*v1.VerifySmartContractWalletSignaturesResponse, error) {
	out := new(v1.VerifySmartContractWalletSignaturesResponse)
	err := c.cc.Invoke(ctx, ValidationApi_VerifySmartContractWalletSignatures_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ValidationApiServer is the server API for ValidationApi service.
// All implementations must embed UnimplementedValidationApiServer
// for forward compatibility
type ValidationApiServer interface {
	// Validates and parses a group message and returns relevant details
	ValidateGroupMessages(context.Context, *ValidateGroupMessagesRequest) (*ValidateGroupMessagesResponse, error)
	// Gets the final association state for a batch of identity updates
	GetAssociationState(context.Context, *GetAssociationStateRequest) (*GetAssociationStateResponse, error)
	// Validates InboxID key packages and returns credential information for them, without checking
	// whether an InboxId <> InstallationPublicKey pair is really valid.
	ValidateInboxIdKeyPackages(context.Context, *ValidateKeyPackagesRequest) (*ValidateInboxIdKeyPackagesResponse, error)
	// Verifies smart contracts
	// This request is proxied from the node, so we'll reuse those messages.
	VerifySmartContractWalletSignatures(context.Context, *v1.VerifySmartContractWalletSignaturesRequest) (*v1.VerifySmartContractWalletSignaturesResponse, error)
	mustEmbedUnimplementedValidationApiServer()
}

// UnimplementedValidationApiServer must be embedded to have forward compatible implementations.
type UnimplementedValidationApiServer struct {
}

func (UnimplementedValidationApiServer) ValidateGroupMessages(context.Context, *ValidateGroupMessagesRequest) (*ValidateGroupMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateGroupMessages not implemented")
}
func (UnimplementedValidationApiServer) GetAssociationState(context.Context, *GetAssociationStateRequest) (*GetAssociationStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAssociationState not implemented")
}
func (UnimplementedValidationApiServer) ValidateInboxIdKeyPackages(context.Context, *ValidateKeyPackagesRequest) (*ValidateInboxIdKeyPackagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateInboxIdKeyPackages not implemented")
}
func (UnimplementedValidationApiServer) VerifySmartContractWalletSignatures(context.Context, *v1.VerifySmartContractWalletSignaturesRequest) (*v1.VerifySmartContractWalletSignaturesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifySmartContractWalletSignatures not implemented")
}
func (UnimplementedValidationApiServer) mustEmbedUnimplementedValidationApiServer() {}

// UnsafeValidationApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ValidationApiServer will
// result in compilation errors.
type UnsafeValidationApiServer interface {
	mustEmbedUnimplementedValidationApiServer()
}

func RegisterValidationApiServer(s grpc.ServiceRegistrar, srv ValidationApiServer) {
	s.RegisterService(&ValidationApi_ServiceDesc, srv)
}

func _ValidationApi_ValidateGroupMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateGroupMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidationApiServer).ValidateGroupMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ValidationApi_ValidateGroupMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidationApiServer).ValidateGroupMessages(ctx, req.(*ValidateGroupMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ValidationApi_GetAssociationState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAssociationStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidationApiServer).GetAssociationState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ValidationApi_GetAssociationState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidationApiServer).GetAssociationState(ctx, req.(*GetAssociationStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ValidationApi_ValidateInboxIdKeyPackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateKeyPackagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidationApiServer).ValidateInboxIdKeyPackages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ValidationApi_ValidateInboxIdKeyPackages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidationApiServer).ValidateInboxIdKeyPackages(ctx, req.(*ValidateKeyPackagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ValidationApi_VerifySmartContractWalletSignatures_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.VerifySmartContractWalletSignaturesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidationApiServer).VerifySmartContractWalletSignatures(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ValidationApi_VerifySmartContractWalletSignatures_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidationApiServer).VerifySmartContractWalletSignatures(ctx, req.(*v1.VerifySmartContractWalletSignaturesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ValidationApi_ServiceDesc is the grpc.ServiceDesc for ValidationApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ValidationApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xmtp.mls_validation.v1.ValidationApi",
	HandlerType: (*ValidationApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateGroupMessages",
			Handler:    _ValidationApi_ValidateGroupMessages_Handler,
		},
		{
			MethodName: "GetAssociationState",
			Handler:    _ValidationApi_GetAssociationState_Handler,
		},
		{
			MethodName: "ValidateInboxIdKeyPackages",
			Handler:    _ValidationApi_ValidateInboxIdKeyPackages_Handler,
		},
		{
			MethodName: "VerifySmartContractWalletSignatures",
			Handler:    _ValidationApi_VerifySmartContractWalletSignatures_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mls_validation/v1/service.proto",
}
