package runner

import (
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func dockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}

func dockerNetworkCreateOptions() network.CreateOptions {
	return network.CreateOptions{
		Driver: "bridge",
	}
}
