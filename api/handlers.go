package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"user_api_with_concurrency/models"
	"user_api_with_concurrency/utils"
)

// Global variables to store users, manage concurrency, and track the next user ID.
var (
	users   = make(map[int]models.User) // Map to store users by their ID.
	usersMu sync.Mutex                  // Mutex to ensure thread-safe access to the users map.
	nextID  = 1                         // Counter to assign unique IDs to new users.
)

// CreateUser handles the creation of a new user.
// It decodes the JSON payload from the request, assigns a unique ID, and stores the user in the map.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Return 400 if the payload is invalid.
		return
	}

	usersMu.Lock()       // Lock the mutex to ensure thread-safe access.
	user.ID = nextID     // Assign the next available ID to the user.
	users[nextID] = user // Add the user to the map.
	nextID++             // Increment the ID counter.
	usersMu.Unlock()     // Unlock the mutex.

	utils.SendUsersToCSV(users) // Export the updated user list to a CSV file.

	w.WriteHeader(http.StatusCreated) // Return 201 (Created) status code.
	json.NewEncoder(w).Encode(user)   // Return the created user as JSON.
}

// GetUsers retrieves all users from the map and returns them as a JSON array.
func GetUsers(w http.ResponseWriter, r *http.Request) {
	usersMu.Lock()         // Lock the mutex to ensure thread-safe access.
	defer usersMu.Unlock() // Ensure the mutex is unlocked when the function exits.

	// Convert the map of users to a slice.
	userList := make([]models.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	w.WriteHeader(http.StatusOK)        // Return 200 (OK) status code.
	json.NewEncoder(w).Encode(userList) // Return the list of users as JSON.
}

// GetUserByID retrieves a specific user by their ID.
// It extracts the ID from the URL, checks if the user exists, and returns the user as JSON.
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Note: r.PathValue("id") only works when the request is made to a router that supports route parameters, such as httprouter or gorilla/mux.
	// However, we are using the standard net/http package, which does not support route parameters directly.
	// To solve this, the ID is being extracted from the URL.
	id, ok := extractUserID(r, w) // Extract the user ID from the URL.
	if !ok {
		return // If the ID is invalid, return an error response.
	}

	usersMu.Lock()         // Lock the mutex to ensure thread-safe access.
	defer usersMu.Unlock() // Ensure the mutex is unlocked when the function exits.

	user, exists := users[id] // Retrieve the user from the map.
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound) // Return 404 if the user doesn't exist.
		return
	}

	w.WriteHeader(http.StatusOK)    // Return 200 (OK) status code.
	json.NewEncoder(w).Encode(user) // Return the user as JSON.
}

// UpdateUser updates an existing user by their ID.
// It extracts the ID from the URL, decodes the updated user data, and updates the user in the map.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := extractUserID(r, w) // Extract the user ID from the URL.
	if !ok {
		return // If the ID is invalid, return an error response.
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Return 400 if the payload is invalid.
		return
	}

	usersMu.Lock()         // Lock the mutex to ensure thread-safe access.
	defer usersMu.Unlock() // Ensure the mutex is unlocked when the function exits.

	user, exists := users[id] // Retrieve the user from the map.
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound) // Return 404 if the user doesn't exist.
		return
	}

	// Update the user's fields.
	user.Name = updatedUser.Name
	user.Age = updatedUser.Age
	user.Email = updatedUser.Email
	users[id] = user // Save the updated user back to the map.

	utils.SendUsersToCSV(users) // Export the updated user list to a CSV file.

	w.WriteHeader(http.StatusOK)    // Return 200 (OK) status code.
	json.NewEncoder(w).Encode(user) // Return the updated user as JSON.
}

// DeleteUser deletes a user by their ID.
// It extracts the ID from the URL, checks if the user exists, and removes them from the map.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := extractUserID(r, w) // Extract the user ID from the URL.
	if !ok {
		return // If the ID is invalid, return an error response.
	}

	usersMu.Lock()         // Lock the mutex to ensure thread-safe access.
	defer usersMu.Unlock() // Ensure the mutex is unlocked when the function exits.

	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound) // Return 404 if the user doesn't exist.
		return
	}

	delete(users, id) // Delete the user from the map.

	utils.SendUsersToCSV(users) // Export the updated user list to a CSV file.

	w.WriteHeader(http.StatusNoContent) // Return 204 (No Content) status code.
}

// extractUserID extracts the user ID from the URL path.
// It validates the ID and returns it as an integer. If the ID is invalid, it returns an error response.
func extractUserID(r *http.Request, w http.ResponseWriter) (int, bool) {
	parts := strings.Split(r.URL.Path, "/") // Split the URL path by "/".
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest) // Return 400 if the URL is invalid.
		return 0, false
	}

	id, err := strconv.Atoi(parts[2]) // Convert the ID from string to integer.
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest) // Return 400 if the ID is invalid.
		return 0, false
	}

	return id, true // Return the valid ID.
}
