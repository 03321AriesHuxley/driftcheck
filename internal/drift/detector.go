package drift

import (
	"fmt"

	"github.com/user/driftcheck/internal/compose"
	"github.com/user/driftcheck/internal/docker"
)

// DriftResult holds the comparison result for a single service.
type DriftResult struct {
	ServiceName string
	Drifts      []string
	Clean       bool
}

// Detector compares running containers against Compose definitions.
type Detector struct {
	inspector *docker.Inspector
}

// NewDetector creates a new Detector with the provided Inspector.
func NewDetector(inspector *docker.Inspector) *Detector {
	return &Detector{inspector: inspector}
}

// Check compares a single service from the Compose file against the running container.
func (d *Detector) Check(service compose.Service, containerName string) (*DriftResult, error) {
	info, err := d.inspector.InspectContainer(containerName)
	if err != nil {
		return nil, fmt.Errorf("inspect container %q: %w", containerName, err)
	}

	result := &DriftResult{
		ServiceName: service.Name,
	}

	// Check image
	if info.Image != service.Image {
		result.Drifts = append(result.Drifts,
			fmt.Sprintf("image mismatch: running=%q compose=%q", info.Image, service.Image))
	}

	// Check environment variables
	runningEnv := toEnvMap(info.Env)
	for _, kv := range service.Environment {
		key, val := splitEnv(kv)
		if rv, ok := runningEnv[key]; !ok {
			result.Drifts = append(result.Drifts, fmt.Sprintf("missing env var: %s", key))
		} else if rv != val {
			result.Drifts = append(result.Drifts,
				fmt.Sprintf("env var %s mismatch: running=%q compose=%q", key, rv, val))
		}
	}

	result.Clean = len(result.Drifts) == 0
	return result, nil
}
