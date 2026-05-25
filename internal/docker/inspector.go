// Package docker provides utilities for inspecting running Docker containers
// and extracting their configuration for drift comparison.
package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerInfo holds the relevant configuration extracted from a running container.
type ContainerInfo struct {
	// Name is the container name (without leading slash).
	Name string
	// Image is the full image reference used to start the container.
	Image string
	// Env is the list of environment variables in KEY=VALUE format.
	Env []string
	// Ports maps container port/protocol to host bindings.
	Ports map[string][]PortBinding
	// Volumes is the list of volume mounts.
	Volumes []VolumeMount
	// Labels contains the container labels.
	Labels map[string]string
}

// PortBinding represents a single host binding for a container port.
type PortBinding struct {
	HostIP   string
	HostPort string
}

// VolumeMount represents a single volume or bind-mount on a container.
type VolumeMount struct {
	Source      string
	Destination string
	Mode        string
}

// Inspector wraps a Docker client and provides high-level inspection helpers.
type Inspector struct {
	cli *client.Client
}

// NewInspector creates an Inspector using the default Docker environment
// (DOCKER_HOST, DOCKER_CERT_PATH, etc.).
func NewInspector() (*Inspector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker: failed to create client: %w", err)
	}
	return &Inspector{cli: cli}, nil
}

// Close releases the underlying Docker client resources.
func (i *Inspector) Close() error {
	return i.cli.Close()
}

// GetContainerByName returns the ContainerInfo for the first running container
// whose name matches the given service name (exact match after stripping the
// leading slash Docker adds internally).
func (i *Inspector) GetContainerByName(ctx context.Context, name string) (*ContainerInfo, error) {
	f := filters.NewArgs(filters.Arg("name", name))
	containers, err := i.cli.ContainerList(ctx, types.ContainerListOptions{
		All:     false,
		Filters: f,
	})
	if err != nil {
		return nil, fmt.Errorf("docker: list containers: %w", err)
	}

	for _, c := range containers {
		for _, n := range c.Names {
			// Docker prefixes names with "/"
			if n == "/"+name || n == name {
				return i.inspect(ctx, c.ID)
			}
		}
	}

	return nil, fmt.Errorf("docker: container %q not found or not running", name)
}

// inspect performs a low-level ContainerInspect call and maps the result to
// ContainerInfo.
func (i *Inspector) inspect(ctx context.Context, id string) (*ContainerInfo, error) {
	data, err := i.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("docker: inspect %s: %w", id, err)
	}

	info := &ContainerInfo{
		Name:   trimLeadingSlash(data.Name),
		Image:  data.Config.Image,
		Env:    data.Config.Env,
		Labels: data.Config.Labels,
		Ports:  make(map[string][]PortBinding),
	}

	// Map port bindings
	for port, bindings := range data.HostConfig.PortBindings {
		key := string(port)
		for _, b := range bindings {
			info.Ports[key] = append(info.Ports[key], PortBinding{
				HostIP:   b.HostIP,
				HostPort: b.HostPort,
			})
		}
	}

	// Map volume mounts
	for _, m := range data.Mounts {
		info.Volumes = append(info.Volumes, VolumeMount{
			Source:      m.Source,
			Destination: m.Destination,
			Mode:        m.Mode,
		})
	}

	return info, nil
}

func trimLeadingSlash(s string) string {
	if len(s) > 0 && s[0] == '/' {
		return s[1:]
	}
	return s
}
