package main

import (
	"flag"
	"fmt"
	"os"
	"user_api_with_concurrency/services"
)

// printUsage displays the usage instructions for the CLI.
// It provides information about available commands and how to use them.
func printUsage() {
	fmt.Println("Usage: cli <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  fetch-additional-info  Fetch additional information for a user")
	fmt.Println()
	fmt.Println("Use './cli <command> --help' for more information on a specific command.")
}

// main is the entry point of the CLI application.
// It parses command-line arguments and executes the appropriate command.
func main() {
	// Check if at least one command is provided.
	if len(os.Args) < 2 {
		printUsage() // Display usage instructions if no command is provided.
		return
	}

	// Switch statement to handle different commands.
	switch os.Args[1] {
	case "fetch-additional-info":
		// Create a new flag set for the "fetch-additional-info" command.
		fetchCmd := flag.NewFlagSet("fetch-additional-info", flag.ExitOnError)
		// Define a flag for the user ID.
		userID := fetchCmd.Int("id", 0, "User ID to fetch information for")
		// Customize the usage message for this command.
		fetchCmd.Usage = func() {
			fmt.Println("Usage: cli fetch-additional-info -id <user_id>")
			fmt.Println("Options:")
			fetchCmd.PrintDefaults()
		}

		// Display help if the "--help" flag is provided.
		if len(os.Args) > 2 && os.Args[2] == "--help" {
			fetchCmd.Usage()
			return
		}

		// Parse the command-line arguments for this command.
		fetchCmd.Parse(os.Args[2:])

		// Validate that the user ID is provided.
		if *userID == 0 {
			fmt.Println("Error: User ID must be specified using -id.")
			return
		}

		// Fetch additional information for the specified user ID.
		users := services.FetchAllUsersInfo([]int{*userID})
		// Print the fetched user information.
		fmt.Println(users)

	default:
		// Handle invalid commands.
		fmt.Println("Error: Invalid command.")
		printUsage() // Display usage instructions for invalid commands.
	}
}
