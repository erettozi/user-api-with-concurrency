package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"user_api_with_concurrency/models"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ProcessAndWriteToCSV processes users from a channel and writes them to a CSV file.
// It filters users to include only those aged 18 or older and formats their names in title case.
func ProcessAndWriteToCSV(userChan <-chan models.User, filename string) error {
	// Create the CSV file.
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed when done.

	// Create a CSV writer.
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all buffered data is written to the file.

	// Create a title caser for formatting names.
	titleCaser := cases.Title(language.English)

	// Write the CSV header.
	if err := writer.Write([]string{"ID", "Name", "Age", "Email"}); err != nil {
		return err
	}

	// Process users from the channel and write them to the CSV file.
	for user := range userChan {
		if user.Age >= 18 { // Only include users aged 18 or older.
			record := []string{
				strconv.Itoa(user.ID),        // Convert ID to string.
				titleCaser.String(user.Name), // Format name in title case.
				strconv.Itoa(user.Age),       // Convert age to string.
				user.Email,                   // Include email.
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}

	// Check for errors during writing.
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

// SendUsersToCSV writes a list or map of users to a CSV file.
// It supports both slices and maps of users and allows specifying a custom filename.
func SendUsersToCSV(users any, filename ...string) {
	var userSlice []models.User

	// Convert the input to a slice of users.
	switch u := users.(type) {
	case []models.User:
		userSlice = u
	case map[int]models.User:
		userSlice = make([]models.User, 0, len(u))
		for _, user := range u {
			userSlice = append(userSlice, user)
		}
	default:
		log.Println("Invalid type for users. Expected []models.User or map[int]models.User")
		return
	}

	// Determine the filename.
	file := "users.csv"
	if len(filename) > 0 && filename[0] != "" {
		file = filename[0]
	}

	// Use a temporary directory for test environments.
	if os.Getenv("ENV") == "test" {
		file = "/tmp/" + file
	}

	// Get the full file path by joining the project root directory with the filename.
	filePath := filepath.Join(GetProjectRoot(), file)

	// Create a channel to send users.
	userChan := make(chan models.User)

	// Start a goroutine to send users to the channel.
	go func() {
		defer close(userChan) // Close the channel when done.
		for _, u := range userSlice {
			userChan <- u
		}
	}()

	// Start a goroutine to process and write users to the CSV file.
	go func() {
		if err := ProcessAndWriteToCSV(userChan, filePath); err != nil {
			log.Println("Failed to process and write CSV:", err)
		}
	}()
}

// GetProjectRoot returns the root directory of the project.
// It uses the runtime.Caller function to determine the location of the current file.
func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)      // Get the path of the current file.
	return filepath.Dir(filepath.Dir(filename)) // Return the project root directory.
}
