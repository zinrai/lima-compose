package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executor handles limactl command execution
type Executor struct{}

// NewExecutor creates a new Executor instance
func NewExecutor() *Executor {
	return &Executor{}
}

// CreateInstance creates a Lima VM instance
func (e *Executor) CreateInstance(name string, instance Instance) error {
	// Build limactl create command with --yes for non-interactive mode
	args := []string{"create", "--name", name, "--yes"}

	// Parse and append args from YAML
	if instance.Args != "" {
		argsList := e.parseArgs(instance.Args)
		args = append(args, argsList...)
	}

	// Append template at the end
	args = append(args, instance.Template)

	// Run the command and let limactl handle existing instance case
	return e.runLimactl(args)
}

// DestroyInstance destroys a Lima VM instance
func (e *Executor) DestroyInstance(name string) error {
	args := []string{"delete", name}
	// Let limactl handle the case where instance doesn't exist
	return e.runLimactl(args)
}

// StartInstance starts a Lima VM instance
func (e *Executor) StartInstance(name string) error {
	args := []string{"start", name}
	return e.runLimactl(args)
}

// StopInstance stops a Lima VM instance
func (e *Executor) StopInstance(name string) error {
	args := []string{"stop", name}
	return e.runLimactl(args)
}

// parseArgs parses the multi-line args string into a slice
func (e *Executor) parseArgs(argsStr string) []string {
	var args []string

	// Split by lines and process each line
	lines := strings.Split(argsStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split by first space to handle flags with values
		parts := strings.SplitN(line, " ", 2)
		for _, part := range parts {
			if part != "" {
				args = append(args, part)
			}
		}
	}

	return args
}

// runLimactl executes a limactl command
func (e *Executor) runLimactl(args []string) error {
	// Always show what command is being executed
	fmt.Printf("=> limactl %s\n", strings.Join(args, " "))

	cmd := exec.Command("limactl", args...)

	// Pass through stdout and stderr directly
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("limactl command failed: %w", err)
	}

	return nil
}
