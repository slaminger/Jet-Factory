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
func PreChroot(mount string) error {
	err := Copy("/usr/bin/qemu-"+buildarch+"-static", mount+"/usr/bin")
	err = SpawnContainer(
		[]string{
			"mount", "--bind",
			mount, mount,

			"&&", "mount", "--bind",
			mount + "/bootloader",
			mount + "/boot",
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
func PostChroot(mounted string) error {
	err := SpawnContainer(
		[]string{
			"rm", mounted + "/usr/bin/qemu-" + buildarch + "-static",
			"&&", "umount", mounted,
			"&&", "mount", mounted + "/boot",
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
func SpawnContainer(cmd, env []string, volume string) error {
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
				Source: volume,
				Target: "/root/" + distribution.Name,
			},
			{
				Type:   mount.TypeVolume,
				Source: client.DefaultDockerHost,
				Target: client.DefaultDockerHost,
			},
		},
	}, nil, distribution.Name)
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
