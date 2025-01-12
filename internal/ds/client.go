package ds

import (
	"fmt"
	"os"

	"github.com/henryppercy/accountability-api/internal/http"
)


var token = os.Getenv("DS_TOKEN")

func FetchUserData() (UserData, error) {
	url := os.Getenv("DS_USER_URL")

	var userResponse UserDataResponse
	err := http.Fetch("GET", url, token, nil, &userResponse, nil)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to fetch user meta times: %w", err)
	}

	return userResponse.User, nil
}

func FetchDailyTimes() (DailyTimes, error) {
	url := os.Getenv("DS_WATCH_TIME_URL")
	var timesResponse DailyTimes
	err := http.Fetch("GET", url, token, nil, &timesResponse, nil)
	if err != nil {
		return DailyTimes{}, fmt.Errorf("failed to fetch all user times: %w", err)
	}

	return timesResponse, nil
}
