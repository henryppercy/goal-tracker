package hevy

import (
	"fmt"
	"os"

	"github.com/henryppercy/accountability-api/internal/http"
)

func FetchWorkouts() (Workouts, error) {
	baseURL := os.Getenv("HEVY_WORKOUTS_URL")

	params := map[string]string{
		"page":   "1",
		"pageSize":  "10",
	}

	headers := map[string]string{
		"api-key": os.Getenv("HEVY_TOKEN"),
	}

	fullURL, _ := http.AddQueryParams(baseURL, params)

	var workoutResponse WorkoutResponse
	err := http.Fetch("GET", fullURL, "", nil, &workoutResponse, headers)
	if err != nil {
		return Workouts{}, fmt.Errorf("failed to fetch user workouts: %w", err)
	}

	return workoutResponse.Workouts, nil
}
