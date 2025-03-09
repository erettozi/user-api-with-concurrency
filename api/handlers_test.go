package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_api_with_concurrency/models"
)

// TestCreateUser tests the CreateUser function.
// It sends a POST request with a JSON payload to create a new user and verifies the response.
func TestCreateUser(t *testing.T) {
	payload := []byte(`{"name":"Erick Rettozi","age":48,"email":"erettozi@tolkien.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	CreateUser(w, req)

	// Check if the status code is 201 (Created).
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Decode the response body into a User struct.
	var user models.User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the user data matches the expected values.
	if user.Name != "Erick Rettozi" || user.Age != 48 || user.Email != "erettozi@tolkien.com" {
		t.Errorf("Unexpected user data: %+v", user)
	}
}

// TestGetUsers tests the GetUsers function.
// It sends a GET request to retrieve all users and verifies the response.
func TestGetUsers(t *testing.T) {
	// Pre-populate the users map with a test user.
	users[1] = models.User{ID: 1, Name: "Erick Rettozi", Age: 48, Email: "erettozi@tolkien.com"}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	GetUsers(w, req)

	// Check if the status code is 200 (OK).
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Decode the response body into a slice of User structs.
	var userList []models.User
	if err := json.NewDecoder(w.Body).Decode(&userList); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify that the number of users returned is correct.
	if len(userList) != 1 {
		t.Errorf("Expected 1 user, got %d", len(userList))
	}
}

// TestGetUserByID tests the GetUserByID function.
// It sends a GET request to retrieve a specific user by ID and verifies the response.
func TestGetUserByID(t *testing.T) {
	// Pre-populate the users map with a test user.
	users[1] = models.User{ID: 1, Name: "Erick Rettozi", Age: 48, Email: "erettozi@tolkien.com"}

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	GetUserByID(w, req)

	// Check if the status code is 200 (OK).
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Decode the response body into a User struct.
	var user models.User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify that the user ID matches the expected value.
	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}
}

// TestUpdateUser tests the UpdateUser function.
// It sends a PUT request to update a specific user by ID and verifies the response.
func TestUpdateUser(t *testing.T) {
	// Pre-populate the users map with a test user.
	users[1] = models.User{ID: 1, Name: "Erick Rettozi", Age: 48, Email: "erettozi@tolkien.com"}

	payload := []byte(`{"name":"Aragorn Elessar","age":37,"email":"aragorn@tolkien.com"}`)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UpdateUser(w, req)

	// Check if the status code is 200 (OK).
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Decode the response body into a User struct.
	var user models.User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify that the user data has been updated correctly.
	if user.Name != "Aragorn Elessar" || user.Age != 37 || user.Email != "aragorn@tolkien.com" {
		t.Errorf("Unexpected user data after update: %+v", user)
	}
}

// TestDeleteUser tests the DeleteUser function.
// It sends a DELETE request to delete a specific user by ID and verifies the response.
func TestDeleteUser(t *testing.T) {
	// Pre-populate the users map with a test user.
	users[1] = models.User{ID: 1, Name: "Erick Rettozi", Age: 48, Email: "erettozi@tolkien.com"}

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	DeleteUser(w, req)

	// Check if the status code is 204 (No Content).
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify that the user has been deleted from the map.
	if _, exists := users[1]; exists {
		t.Error("User was not deleted")
	}
}
