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

// BinfmtSupport : Adds Binfmt support for cross-architecture chroot; returns nil on success; returns err otherwise;
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

// SpawnContainer : Spawns a container based on dockerImageName; returns nil on success; returns err otherwise;
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

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	config := &container.HostConfig{
		Privileged:  true,
		VolumesFrom: []string{hostname},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        dockerImageName,
		Entrypoint:   []string{"/bin/bash"},
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
