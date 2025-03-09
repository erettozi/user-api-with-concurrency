package models

// User represents a user entity in the application.
// It defines the structure of a user, including their ID, name, age, and email.
type User struct {
	ID    int    `json:"id"`    // Unique identifier for the user.
	Name  string `json:"name"`  // Full name of the user.
	Age   int    `json:"age"`   // Age of the user.
	Email string `json:"email"` // Email address of the user.
}
