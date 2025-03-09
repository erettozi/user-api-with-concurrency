package main

import (
	"os"
	"os/exec"
	"testing"
)

// TestCLI tests the command-line interface (CLI) of the application.
// It runs the main program with specific arguments and checks if the command executes successfully.
func TestCLI(t *testing.T) {
	// Create a command to run the main program with the "fetch-additional-info" subcommand and an ID flag.
	cmd := exec.Command("go", "run", "main.go", "fetch-additional-info", "-id=1")

	// Redirect the command's standard output and error to the test's standard output and error.
	// This allows the test to display the output of the command in real-time.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and check for errors.
	if err := cmd.Run(); err != nil {
		// If the command fails, report the error and mark the test as failed.
		t.Errorf("CLI command failed: %v", err)
	}
}
