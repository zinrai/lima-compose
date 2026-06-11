package main

import (
	"fmt"
	"os"
)

const defaultComposeFile = "lima-compose.yaml"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "create":
		runCreate(args)
	case "delete":
		runDelete(args)
	case "start":
		runStart(args)
	case "stop":
		runStop(args)
	case "hosts":
		runHosts(args)
	case "version":
		runVersion()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("lima-compose - Compose for Lima VMs")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  lima-compose <command> [compose-file]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create    Create all instances defined in the compose file")
	fmt.Println("  delete    Delete all instances defined in the compose file")
	fmt.Println("  start     Start all instances defined in the compose file")
	fmt.Println("  stop      Stop all instances defined in the compose file")
	fmt.Println("  hosts     Show /etc/hosts entries for all running instances")
	fmt.Println("  version   Show version information")
	fmt.Println("  help      Show this help message")
	fmt.Println()
	fmt.Println("If compose-file is not specified, 'lima-compose.yaml' will be used.")
}

func getComposeFile(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	// Try lima-compose.yaml first, then lima-compose.yml
	if _, err := os.Stat(defaultComposeFile); err == nil {
		return defaultComposeFile
	}

	altFile := "lima-compose.yml"
	if _, err := os.Stat(altFile); err == nil {
		return altFile
	}

	return defaultComposeFile // Return default even if not exists
}

func runCreate(args []string) {
	file := getComposeFile(args)
	compose, err := LoadCompose(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading compose file: %v\n", err)
		os.Exit(1)
	}

	executor := NewExecutor()
	hasError := false

	for name, instance := range compose.Instances {
		fmt.Printf("Creating instance: %s\n", name)
		if err := executor.CreateInstance(name, instance); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", name, err)
			hasError = true
			// Continue with other instances
		}
	}

	if hasError {
		os.Exit(1)
	}
}

func runDelete(args []string) {
	file := getComposeFile(args)
	compose, err := LoadCompose(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading compose file: %v\n", err)
		os.Exit(1)
	}

	executor := NewExecutor()
	hasError := false

	for name := range compose.Instances {
		fmt.Printf("Deleting instance: %s\n", name)
		if err := executor.DeleteInstance(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting %s: %v\n", name, err)
			hasError = true
			// Continue with other instances
		}
	}

	if hasError {
		os.Exit(1)
	}
}

func runStart(args []string) {
	file := getComposeFile(args)
	compose, err := LoadCompose(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading compose file: %v\n", err)
		os.Exit(1)
	}

	executor := NewExecutor()
	hasError := false

	for name := range compose.Instances {
		fmt.Printf("Starting instance: %s\n", name)
		if err := executor.StartInstance(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting %s: %v\n", name, err)
			hasError = true
		}
	}

	if hasError {
		os.Exit(1)
	}
}

func runStop(args []string) {
	file := getComposeFile(args)
	compose, err := LoadCompose(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading compose file: %v\n", err)
		os.Exit(1)
	}

	executor := NewExecutor()
	errors := []error{}

	for name := range compose.Instances {
		fmt.Printf("Stopping instance: %s\n", name)
		if err := executor.StopInstance(name); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		}
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "\nFailed to stop all instances:\n")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "  - %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("\nAll instances stopped successfully.")
}

func runVersion() {
	fmt.Printf("lima-compose %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("date:   %s\n", date)
}

func runHosts(args []string) {
	file := getComposeFile(args)
	compose, err := LoadCompose(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading compose file: %v\n", err)
		os.Exit(1)
	}

	PrintHostsFormat(compose)
}
