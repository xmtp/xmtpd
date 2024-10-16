package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// / Maps from a URL, as defined in https://pkg.go.dev/net/url#URL, to a gRPC target,
// / as defined in https://github.com/grpc/grpc/blob/master/doc/naming.md
func HttpAddressToGrpcTarget(httpAddress string) (string, error) {
	url, err := url.Parse(httpAddress)
	if err != nil {
		return "", err
	}
	if strings.ToLower(url.Hostname()) == "localhost" {
		return fmt.Sprintf("passthrough://localhost/[::]:%s", url.Port()), nil
	}

	// TODO(rich): Handle SSL properly
	return fmt.Sprintf("dns:%s:%s", url.Hostname(), url.Port()), nil
}
