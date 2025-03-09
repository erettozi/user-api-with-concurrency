package tests

import (
	"testing"
	"user_api_with_concurrency/services"
)

// BenchmarkFetchAllUsersInfo benchmarks the performance of the FetchAllUsersInfo function.
// It measures how long it takes to fetch user information for a given list of user IDs.
func BenchmarkFetchAllUsersInfo(b *testing.B) {
	// Define a list of user IDs to test.
	userIDs := []int{1, 2, 3, 4, 5}

	// Run the FetchAllUsersInfo function b.N times to measure its performance.
	for i := 0; i < b.N; i++ {
		services.FetchAllUsersInfo(userIDs)
	}
}
