package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// HTTPAddressToGRPCTarget maps from a URL, as defined in https://pkg.go.dev/net/url#URL, to a gRPC target,
// as defined in https://github.com/grpc/grpc/blob/master/doc/naming.md
func HTTPAddressToGRPCTarget(httpAddress string) (string, bool, error) {
	url, err := url.Parse(httpAddress)
	if err != nil {
		return "", false, err
	}
	var isTLS bool
	switch url.Scheme {
	case "https":
		isTLS = true
	case "http", "":
		isTLS = false
	default:
		return "", false, fmt.Errorf("unknown connection schema %s", url.Scheme)
	}

	if url.Port() != "" {
		return fmt.Sprintf("%s:%s", url.Hostname(), url.Port()), isTLS, nil
	}

	return url.Hostname(), isTLS, nil
}

func GetCredentialsForAddress(
	isTLS bool,
) (credentials.TransportCredentials, error) {
	if !isTLS {
		return insecure.NewCredentials(), nil
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system CA certificates: %v", err)
	}
	if certPool == nil {
		return nil, fmt.Errorf("no system CA certificates available")
	}

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool, // System CA pool
	})

	return creds, nil
}
