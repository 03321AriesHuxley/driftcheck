package compose

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents a single service definition from a Compose file.
type ServiceConfig struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
	Command     string            `yaml:"command"`
	Restart     string            `yaml:"restart"`
}

// ComposeFile represents the top-level structure of a docker-compose.yml file.
type ComposeFile struct {
	Version  string                   `yaml:"version"`
	Services map[string]ServiceConfig `yaml:"services"`
}

// ParseFile reads and parses a Docker Compose YAML file at the given path.
func ParseFile(path string) (*ComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file %q: %w", path, err)
	}

	var cf ComposeFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsing compose file %q: %w", path, err)
	}

	if cf.Services == nil {
		cf.Services = make(map[string]ServiceConfig)
	}

	return &cf, nil
}

// GetService returns the ServiceConfig for the named service, or an error if
// the service is not defined in the Compose file.
func (cf *ComposeFile) GetService(name string) (ServiceConfig, error) {
	svc, ok := cf.Services[name]
	if !ok {
		return ServiceConfig{}, fmt.Errorf("service %q not found in compose file", name)
	}
	return svc, nil
}
