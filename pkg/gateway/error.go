package gateway

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type GatewayServiceError interface {
	error
	Code() codes.Code
	ClientMessage() string
}

type PermissionDeniedError struct {
	msg string
	err error
}

func (e PermissionDeniedError) Error() string {
	if e.err == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e PermissionDeniedError) ClientMessage() string {
	return e.msg
}

func (e PermissionDeniedError) Code() codes.Code {
	return codes.PermissionDenied
}

func NewPermissionDeniedError(msg string, err error) *PermissionDeniedError {
	return &PermissionDeniedError{msg: msg, err: err}
}

type UnauthenticatedError struct {
	msg string
	err error
}

func (e UnauthenticatedError) Error() string {
	if e.err == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e UnauthenticatedError) ClientMessage() string {
	return e.msg
}

func (e UnauthenticatedError) Code() codes.Code {
	return codes.Unauthenticated
}

func NewUnauthenticatedError(msg string, err error) *UnauthenticatedError {
	return &UnauthenticatedError{msg: msg, err: err}
}
