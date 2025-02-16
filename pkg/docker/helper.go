package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/loft-sh/devpod/pkg/devcontainer/config"
	"github.com/loft-sh/devpod/pkg/image"
	"github.com/loft-sh/devpod/pkg/scanner"
	perrors "github.com/pkg/errors"
)

type DockerHelper struct {
	DockerCommand string
}

func (r *DockerHelper) GPUSupportEnabled() (bool, error) {
	out, err := r.buildCmd(context.TODO(), "info", "-f", "{{.Runtimes.nvidia}}").Output()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(out), "nvidia-container-runtime"), nil
}

func (r *DockerHelper) FindDevContainer(labels []string) (*config.ContainerDetails, error) {
	containers, err := r.FindContainer(labels)
	if err != nil {
		return nil, err
	} else if len(containers) == 0 {
		return nil, nil
	}

	containerDetails, err := r.InspectContainers(containers)
	if err != nil {
		return nil, err
	}

	// find matching container
	for _, details := range containerDetails {
		if strings.ToLower(details.State.Status) != "removing" {
			details.State.Status = strings.ToLower(details.State.Status)
			return &details, nil
		}
	}

	return nil, nil
}

func (r *DockerHelper) Stop(ctx context.Context, id string) error {
	out, err := r.buildCmd(ctx, "stop", id).CombinedOutput()
	if err != nil {
		return perrors.Wrapf(err, "%s", string(out))
	}

	return nil
}

func (r *DockerHelper) Remove(ctx context.Context, id string) error {
	out, err := r.buildCmd(ctx, "rm", id).CombinedOutput()
	if err != nil {
		return perrors.Wrapf(err, "%s", string(out))
	}

	return nil
}

func (r *DockerHelper) Run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	return r.RunWithDir(ctx, "", args, stdin, stdout, stderr)
}

func (r *DockerHelper) RunWithDir(ctx context.Context, dir string, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := r.buildCmd(ctx, args...)
	cmd.Dir = dir
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func (r *DockerHelper) StartContainer(ctx context.Context, id string, labels []string) error {
	out, err := r.buildCmd(ctx, "start", id).CombinedOutput()
	if err != nil {
		return perrors.Wrapf(err, "start command: %v", string(out))
	}

	container, err := r.FindDevContainer(labels)
	if err != nil {
		return err
	} else if container == nil {
		return fmt.Errorf("container not found")
	}

	return nil
}

func (r *DockerHelper) InspectImage(imageName string, tryRemote bool) (*config.ImageDetails, error) {
	imageDetails := []*config.ImageDetails{}
	err := r.Inspect([]string{imageName}, "image", &imageDetails)
	if err != nil {
		// try remote?
		if !tryRemote {
			return nil, err
		}

		imageConfig, _, err := image.GetImageConfig(imageName)
		if err != nil {
			return nil, perrors.Wrap(err, "get image config remotely")
		}

		return &config.ImageDetails{
			ID: imageName,
			Config: config.ImageDetailsConfig{
				User:       imageConfig.Config.User,
				Env:        imageConfig.Config.Env,
				Labels:     imageConfig.Config.Labels,
				Entrypoint: imageConfig.Config.Entrypoint,
				Cmd:        imageConfig.Config.Cmd,
			},
		}, nil
	} else if len(imageDetails) == 0 {
		return nil, fmt.Errorf("cannot find image details for %s", imageName)
	}

	return imageDetails[0], nil
}

func (r *DockerHelper) InspectContainers(ids []string) ([]config.ContainerDetails, error) {
	details := []config.ContainerDetails{}
	err := r.Inspect(ids, "container", &details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (r *DockerHelper) Inspect(ids []string, inspectType string, obj interface{}) error {
	args := []string{"inspect", "--type", inspectType}
	args = append(args, ids...)
	out, err := r.buildCmd(context.TODO(), args...).Output()
	if err != nil {
		return perrors.Wrapf(err, "inspect container: %v", string(out))
	}

	err = json.Unmarshal(out, obj)
	if err != nil {
		return perrors.Wrap(err, "parse inspect output")
	}

	return nil
}

func (r *DockerHelper) FindContainer(labels []string) ([]string, error) {
	args := []string{"ps", "-q", "-a"}
	for _, label := range labels {
		args = append(args, "--filter", "label="+label)
	}

	out, err := r.buildCmd(context.TODO(), args...).Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return nil, fmt.Errorf("find container: %s%s", string(out), strings.TrimSpace(string(exitError.Stderr)))
		}

		return nil, fmt.Errorf("find container: %s", strings.TrimSpace(string(out)))
	}

	arr := []string{}
	scan := scanner.NewScanner(bytes.NewReader(out))
	for scan.Scan() {
		arr = append(arr, strings.TrimSpace(scan.Text()))
	}

	return arr, nil
}

func (r *DockerHelper) buildCmd(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, r.DockerCommand, args...)
}
