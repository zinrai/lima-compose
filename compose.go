package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Compose represents the lima-compose.yaml structure
type Compose struct {
	Instances map[string]Instance `yaml:"instances"`
}

// Instance represents a single VM instance configuration
type Instance struct {
	Template string `yaml:"template"`
	Args     string `yaml:"args"`
}

// LoadCompose loads and parses a compose YAML file
func LoadCompose(filename string) (*Compose, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return parseYAML(data)
}

// parseYAML parses YAML data into Compose struct
func parseYAML(data []byte) (*Compose, error) {
	var compose Compose

	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if compose.Instances == nil || len(compose.Instances) == 0 {
		return nil, fmt.Errorf("no instances defined in compose file")
	}

	// Validate each instance
	for name, instance := range compose.Instances {
		if instance.Template == "" {
			return nil, fmt.Errorf("instance %s: template is required", name)
		}
	}

	return &compose, nil
}

// GetInstanceNames returns a sorted list of instance names
func (c *Compose) GetInstanceNames() []string {
	names := make([]string, 0, len(c.Instances))
	for name := range c.Instances {
		names = append(names, name)
	}
	return names
}
