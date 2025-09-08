package main

import (
	"reflect"
	"testing"
)

// Core functionality test: Complex argument parsing
// If this breaks, lima VMs won't be created with correct parameters
func TestParseArgs_ComplexArguments(t *testing.T) {
	input := `      --cpus 2
      --memory 4
      --set '.env.NODE_TYPE="server"'
      --mount ~/projects:/projects:w`

	executor := NewExecutor()
	result := executor.parseArgs(input)

	expected := []string{
		"--cpus", "2",
		"--memory", "4",
		"--set", `'.env.NODE_TYPE="server"'`,
		"--mount", "~/projects:/projects:w",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseArgs() failed to handle complex arguments\ngot:  %v\nwant: %v", result, expected)
	}
}
