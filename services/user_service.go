package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"user_api_with_concurrency/models"
	"user_api_with_concurrency/utils"
)

const MaxConcurrentFetches = 5 // Global limit of competition

var (
	externalAPIURL string // Stores the URL of the external API.
)

// init initializes the externalAPIURL variable.
// It reads the URL from the environment variable EXTERNAL_API_URL.
// If the environment variable is not set, it defaults to "http://localhost:3000".
func init() {
	externalAPIURL = os.Getenv("EXTERNAL_API_URL")
	if externalAPIURL == "" {
		externalAPIURL = "http://localhost:3000"
	}
}

// FetchAdditionalInfo fetches additional information for a specific user from the external API.
// It is designed to be run as a goroutine and uses a semaphore to limit concurrency.
func FetchAdditionalInfo(userID int, wg *sync.WaitGroup, results chan<- models.User, semaphore chan struct{}) {
	defer wg.Done()                // Notify the WaitGroup that this goroutine is done.
	defer func() { <-semaphore }() // Release the semaphore slot when done.

	// Make an HTTP GET request to the external API to fetch user information.
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", externalAPIURL, userID))
	if err != nil {
		fmt.Printf("Error fetching user info for user %d: %v\n", userID, err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading.

	// Check if the response status code is not 200 (OK).
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching user info for user %d: received status code %d\n", userID, resp.StatusCode)
		return
	}

	// Decode the JSON response into a User struct.
	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Printf("Error decoding response for user %d: %v\n", userID, err)
		return
	}

	// Send the fetched user data to the results channel.
	results <- user
}

// FetchAllUsersInfo fetches additional information for multiple users concurrently.
// It uses a semaphore to limit the number of concurrent goroutines.
func FetchAllUsersInfo(userIDs []int) []models.User {
	var wg sync.WaitGroup                                  // WaitGroup to wait for all goroutines to finish.
	results := make(chan models.User, len(userIDs))        // Buffered channel to store fetched user data.
	semaphore := make(chan struct{}, MaxConcurrentFetches) // Semaphore to limit concurrency based on the global constant.

	// Iterate over the list of user IDs and start a goroutine for each.
	for _, id := range userIDs {
		wg.Add(1)                                           // Increment the WaitGroup counter.
		semaphore <- struct{}{}                             // Acquire a semaphore slot.
		go FetchAdditionalInfo(id, &wg, results, semaphore) // Start the goroutine.
	}

	wg.Wait()      // Wait for all goroutines to finish.
	close(results) // Close the results channel to signal that no more data will be sent.

	// Collect the fetched user data from the results channel.
	var users []models.User
	for user := range results {
		users = append(users, user)
	}

	// Export the fetched user data to a CSV file.
	utils.SendUsersToCSV(users)

	return users // Return the list of fetched users.
}
