package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"user_api_with_concurrency/models"
)

// TestFetchAllUsersInfo tests the FetchAllUsersInfo function.
// It simulates an external API server and verifies that the function correctly fetches and processes user data.
func TestFetchAllUsersInfo(t *testing.T) {
	// Simulate an HTTP server to mock responses from the external API.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int
		fmt.Sscanf(r.URL.Path, "/users/%d", &id) // Extract the user ID from the URL path.

		// Create a mock user with the extracted ID.
		user := models.User{
			ID:   id,
			Name: fmt.Sprintf("User %d", id),
			Age:  20 + id, // Fictional age for testing.
		}

		// Set the response headers and encode the user as JSON.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}))
	defer ts.Close() // Ensure the test server is closed after the test.

	// Replace the external API URL with the test server URL.
	oldURL := externalAPIURL
	externalAPIURL = ts.URL
	defer func() { externalAPIURL = oldURL }() // Restore the original URL after the test.

	// Define a list of user IDs to test.
	userIDs := []int{1, 2, 3, 4, 5}

	// Call the function under test.
	users := FetchAllUsersInfo(userIDs)

	// Sort the users by ID to ensure consistent comparison.
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	// Verify that the correct number of users was returned.
	if len(users) != len(userIDs) {
		t.Errorf("Expected %d users, got %d", len(userIDs), len(users))
	}

	// Verify the data for each user.
	for i, user := range users {
		expectedID := userIDs[i]
		expectedName := fmt.Sprintf("User %d", expectedID)
		expectedAge := 20 + expectedID

		if user.ID != expectedID {
			t.Errorf("Expected user ID %d, got %d", expectedID, user.ID)
		}
		if user.Name != expectedName {
			t.Errorf("Expected user name %s, got %s", expectedName, user.Name)
		}
		if user.Age != expectedAge {
			t.Errorf("Expected user age %d, got %d", expectedAge, user.Age)
		}
	}
}

// TestFetchAllUsersInfo_ConcurrencyLimit tests the concurrency limit of the FetchAllUsersInfo function.
// It ensures that the function does not exceed the maximum allowed number of concurrent goroutines.
func TestFetchAllUsersInfo_ConcurrencyLimit(t *testing.T) {
	var (
		mu          sync.Mutex
		maxRoutines int // Tracks the maximum number of concurrent goroutines.
		current     int // Tracks the current number of active goroutines.
	)

	// Simulate an HTTP server to mock responses from the external API.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		current++
		if current > maxRoutines {
			maxRoutines = current // Update the maximum number of concurrent goroutines.
		}
		mu.Unlock()

		defer func() {
			mu.Lock()
			current-- // Decrement the count of active goroutines when done.
			mu.Unlock()
		}()

		// Extract the user ID from the URL path.
		parts := strings.Split(r.URL.Path, "/")
		idStr := parts[2]
		userID, _ := strconv.Atoi(idStr)

		// Create a mock user with the extracted ID.
		user := models.User{
			ID:   userID,
			Name: fmt.Sprintf("User %d", userID),
			Age:  20 + userID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}))
	defer ts.Close() // Ensure the test server is closed after the test.

	// Replace the external API URL with the test server URL.
	oldURL := externalAPIURL
	externalAPIURL = ts.URL
	defer func() { externalAPIURL = oldURL }() // Restore the original URL after the test.

	// Define a list of user IDs to test.
	userIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Call the function under test.
	_ = FetchAllUsersInfo(userIDs)

	// Verify that the maximum number of concurrent goroutines does not exceed the limit (5).
	if maxRoutines > 5 {
		t.Errorf("Expected maximum of 5 concurrent goroutines, got %d", maxRoutines)
	}
}

// TestFetchAllUsersInfo_ErrorHandling tests the error handling of the FetchAllUsersInfo function.
// It ensures that the function handles API errors correctly and returns an empty result.
func TestFetchAllUsersInfo_ErrorHandling(t *testing.T) {
	// Simulate an HTTP server that returns an error response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Simulate a server error.
	}))
	defer ts.Close() // Ensure the test server is closed after the test.

	// Replace the external API URL with the test server URL.
	oldURL := externalAPIURL
	externalAPIURL = ts.URL
	defer func() { externalAPIURL = oldURL }() // Restore the original URL after the test.

	// Define a list of user IDs to test.
	userIDs := []int{1, 2, 3}

	// Call the function under test.
	users := FetchAllUsersInfo(userIDs)

	// Verify that no users are returned due to the simulated errors.
	if len(users) != 0 {
		t.Errorf("Expected 0 users due to errors, got %d", len(users))
	}
}
