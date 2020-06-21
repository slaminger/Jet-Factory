package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// BinfmtSupport :
func BinfmtSupport() error {
	ctx := context.Background()
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, "docker.io/aptman/qus", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "docker.io/aptman/qus",
		Cmd:   []string{"-s", "-- -p"},
	}, &container.HostConfig{
		Privileged: true,
	}, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if _, err := cli.ContainerWait(ctx, resp.ID); err != nil {
		return err
	}
	return nil
}

// SpawnContainer : Spawns a container based on dockerImageName
func SpawnContainer(cmd, env []string) error {
	ctx := context.Background()
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, dockerImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	config := &container.HostConfig{
		Privileged:  true,
		VolumesFrom: []string{"jet"},
		// Mounts: []mount.Mount{
		// 	{
		// 		Type:   mount.TypeBind,
		// 		Source: "/var/run/docker.sock",
		// 		Target: "/var/run/docker.sock",
		// 	},
		// },
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        dockerImageName,
		Cmd:          cmd,
		Env:          env,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, config, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	go func() error {
		reader, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Timestamps: false,
		})
		if err != nil {
			return err
		}
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		return nil
	}()

	if _, err := cli.ContainerWait(ctx, resp.ID); err != nil {
		return err
	}

	return nil
}
