package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // Import for pprof (profiling) support.
	"os"
	"user_api_with_concurrency/api"
)

// main is the entry point of the application.
// It starts a pprof server for profiling, sets up API routes, and starts the HTTP server.
func main() {
	// Start a goroutine to run the pprof server for profiling.
	go func() {
		log.Println("Starting pprof server at http://localhost:6060/debug/pprof")
		// Start the pprof server on localhost:6060.
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Set up the API routes using the SetupRoutes function from the api package.
	api.SetupRoutes()

	// Get the port to listen on from the environment variable or use a default value.
	port := getPort()
	log.Printf("Server started on :%s\n", port)

	// Start the HTTP server and listen for incoming requests.
	// If the server fails to start, log the error and exit.
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

// getPort retrieves the port number from the environment variable PORT.
// If the PORT environment variable is not set, it defaults to "3000".
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if PORT is not set.
	}
	return port
}
