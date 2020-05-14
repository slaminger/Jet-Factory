package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// PreChroot : Copy qemu-aarch64-static binary and mount bind the directories
func PreChroot(mount [2]string) error {
	err := SpawnContainer(
		[]string{
			"cp", "/usr/bin/qemu-" + buildarch + "-static",
			mount[1] + "/usr/bin",

			"&&", "mount", "--bind",
			mount[1], mount[1],

			"&&", "mount", "--bind",
			mount[1] + "/bootloader",
			mount[1] + "/boot",
		},
		nil,
		mount,
	)
	if err != nil {
		return err
	}
	return nil
}

// PostChroot : Remove qemu-aarch64-static binary and unmount the binded directories
func PostChroot(mounted [2]string) error {
	err := SpawnContainer(
		[]string{
			"rm", mounted[1] + "/usr/bin/qemu-" + buildarch + "-static",
			"&&", "umount", mounted[1],
			"&&", "mount", mounted[1] + "/boot",
		},
		nil,
		mounted,
	)
	if err != nil {
		return err
	}
	return nil
}

// SpawnContainer : Spawns a container based on dockerImageName
func SpawnContainer(cmd, env []string, volume [2]string) error {
	ctx := context.Background()
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}

	// cli.ImageBuild(ctx)

	reader, err := cli.ImagePull(ctx, dockerImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: dockerImageName,
		Cmd:   cmd,
		Env:   env,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volume[0],
				Target: volume[1],
			},
		},
	}, nil, baseName)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if _, err := cli.ContainerWait(ctx, resp.ID); err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
