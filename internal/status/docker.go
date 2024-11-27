package status

import "github.com/docker/docker/client"

type Docker struct {
	Version string
}

func FetchDocker() (*Docker, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	return &Docker{Version: apiClient.ClientVersion()}, nil
}
