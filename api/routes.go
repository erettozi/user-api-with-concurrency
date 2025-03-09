package api

import "net/http"

// SetupRoutes configures the HTTP routes for the API.
// It maps specific HTTP methods and URL paths to their corresponding handler functions.
func SetupRoutes() {
	// Register the route for creating a new user.
	// When a POST request is made to "/users", the CreateUser function will handle it.
	http.HandleFunc("POST /users", CreateUser)

	// Register the route for retrieving all users.
	// When a GET request is made to "/users", the GetUsers function will handle it.
	http.HandleFunc("GET /users", GetUsers)

	// Register the route for retrieving a specific user by ID.
	// When a GET request is made to "/users/{id}", the GetUserByID function will handle it.
	// The {id} part is a path parameter that represents the user's ID.
	http.HandleFunc("GET /users/{id}", GetUserByID)

	// Register the route for updating a specific user by ID.
	// When a PUT request is made to "/users/{id}", the UpdateUser function will handle it.
	// The {id} part is a path parameter that represents the user's ID.
	http.HandleFunc("PUT /users/{id}", UpdateUser)

	// Register the route for deleting a specific user by ID.
	// When a DELETE request is made to "/users/{id}", the DeleteUser function will handle it.
	// The {id} part is a path parameter that represents the user's ID.
	http.HandleFunc("DELETE /users/{id}", DeleteUser)
}
