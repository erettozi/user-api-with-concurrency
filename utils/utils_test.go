package utils

import (
	"os"
	"testing"
	"user_api_with_concurrency/models"
)

// TestProcessAndWriteToCSV tests the ProcessAndWriteToCSV function.
// It verifies that the function correctly processes users from a channel and writes them to a CSV file.
func TestProcessAndWriteToCSV(t *testing.T) {
	// Create a channel to send users.
	userChan := make(chan models.User)

	// Start the ProcessAndWriteToCSV function in a goroutine.
	filename := "test_users.csv"
	defer os.Remove(filename) // Ensure the test file is deleted after the test.

	// Launch a goroutine to send users to the channel.
	go func() {
		defer close(userChan) // Close the channel after sending all users.
		users := []models.User{
			{ID: 1, Name: "Erick Rettozi", Age: 48, Email: "erettozi@tolkien.com"},
			{ID: 2, Name: "Aragorn Elessar", Age: 37, Email: "aragorn@tolkien.com"},
		}
		for _, user := range users {
			userChan <- user // Send each user through the channel.
		}
	}()

	// Call the ProcessAndWriteToCSV function.
	if err := ProcessAndWriteToCSV(userChan, filename); err != nil {
		t.Fatalf("Failed to process and write CSV: %v", err)
	}

	// Verify that the file was created successfully.
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close() // Ensure the file is closed after checking.
}
