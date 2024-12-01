package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
)

type Docker struct {
	Version    string
	Containers []DockerContainer
}

type DockerContainer struct {
	Id      string
	Image   string
	Command string
	Created string
	Status  string
	Name    string
	Ports   []DockerContainerPort
}

type DockerContainerPort struct {
	Ip          string
	PrivatePort string
	PublicPort  string
	Type        string
}

func newDockerContainerPort(port types.Port) DockerContainerPort {
	return DockerContainerPort{
		Ip:          port.IP,
		PrivatePort: strconv.FormatUint(uint64(port.PrivatePort), 10),
		PublicPort:  strconv.FormatUint(uint64(port.PublicPort), 10),
		Type:        port.Type,
	}
}

func newDockerContainer(container types.Container) DockerContainer {
	ports := make([]DockerContainerPort, 0)
	for _, port := range container.Ports {
		ports = append(ports, newDockerContainerPort(port))
	}

	return DockerContainer{
		Id:      container.ID,
		Image:   container.Image,
		Command: container.Command,
		Created: humanize.Time(time.Unix(container.Created, 0)),
		Status:  container.Status,
		Name:    container.Names[0],
		Ports:   ports,
	}
}

func FetchDocker(ctx context.Context) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	defer cli.Close()

	containerList, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	containers := make([]DockerContainer, 0)
	for _, container := range containerList {
		containers = append(containers, newDockerContainer(container))
	}

	return &Docker{Version: cli.ClientVersion(), Containers: containers}, nil
}

func CreateContainer(ctx context.Context, img string, version string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	defer cli.Close()

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: fmt.Sprintf("%s:%s", img, version),
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		nil,
		"",
	)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func FindContainer(ctx context.Context, id string) (*DockerContainer, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if c.ID == id {
			container := newDockerContainer(c)
			return &container, nil
		}
	}

	return nil, errors.New("container not found")

}

func StartContainer(ctx context.Context, id string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	defer cli.Close()

	return cli.ContainerStart(ctx, id, container.StartOptions{})
}

func StopContainer(ctx context.Context, id string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	defer cli.Close()

	return cli.ContainerStop(ctx, id, container.StopOptions{})
}

func RemoveContainer(ctx context.Context, id string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	defer cli.Close()

	return cli.ContainerRemove(ctx, id, container.RemoveOptions{})
}

type LogMessage struct {
	Error error
	Text  string
}

func StreamLogs(ctx context.Context, id string, logChannel chan LogMessage) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println(err)
		logChannel <- LogMessage{Error: err, Text: ""}
		return
	}

	defer cli.Close()

	reader, err := cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})

	if err != nil {
		log.Println(err)
		logChannel <- LogMessage{Error: err, Text: ""}
		return
	}

	defer reader.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			logChannel <- LogMessage{Error: err, Text: ""}
			return
		}

		text := string(buffer[:n])
		logChannel <- LogMessage{Error: nil, Text: text}
	}
}
