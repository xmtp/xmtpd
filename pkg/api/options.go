package api

import (
	"net"
	"strconv"

	"github.com/pkg/errors"
)

type Options struct {
	GRPCAddress string `long:"grpc-address" description:"API gRPC listen address" default:"0.0.0.0"`
	GRPCPort    uint   `long:"grpc-port" description:"API gRPC listen port" default:"5000"`
	HTTPAddress string `long:"http-address" description:"API HTTP listen address" default:"0.0.0.0"`
	HTTPPort    uint   `long:"http-port" description:"API HTTP listen port" default:"5001"`
	MaxMsgSize  int    `long:"max-msg-size" description:"Max message size in bytes (default 10Mb)" default:"1250000"`
}

func (opts *Options) validate() error {
	if err := validateAddr(opts.HTTPAddress, opts.HTTPPort); err != nil {
		return errors.Wrap(err, "Invalid HTTP Address")
	}
	if err := validateAddr(opts.GRPCAddress, opts.GRPCPort); err != nil {
		return errors.Wrap(err, "Invalid GRPC Address")
	}
	return nil
}

func validateAddr(addr string, port uint) error {
	_, err := net.ResolveTCPAddr("tcp", hostPortAddr(addr, port))
	return err
}

func hostPortAddr(addr string, port uint) string {
	return net.JoinHostPort(addr, strconv.Itoa(int(port)))
}
