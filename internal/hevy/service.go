package hevy

import (
	"database/sql"
	"fmt"
	"time"
)

func GetHevyData(db *sql.DB) (HevyData, error) {
	goal := 2
	totalWorkouts, err := GetWorkouts(db)
	if err != nil {
		return HevyData{}, fmt.Errorf("failed to fetch hevy data: %w", err)
	}

	weekWorkouts := getWeekWorkouts(totalWorkouts)
	todaysWorkout := getTodayWorkout(totalWorkouts)
	yearWorkouts := getYearWorkouts(totalWorkouts)

	allTime := AllTime{
		TotalWorkouts: getTotalWorkouts(totalWorkouts),
		TimeTotal:     getTotalWorkoutTime(totalWorkouts),
		WeeksGoalMet:  getSumGoalMet(goal, totalWorkouts),
	}

	year := Year{
		TotalWorkouts: getTotalWorkouts(yearWorkouts),
		TimeTotal:     getTotalWorkoutTime(yearWorkouts),
		WeeksGoalMet:  getSumGoalMet(goal, yearWorkouts),
	}

	week := Week{
		TotalWorkouts: getTotalWorkouts(weekWorkouts),
		TimeTotal:     getTotalWorkoutTime(weekWorkouts),
		GoalMet:       getGoalMet(goal, weekWorkouts),
	}

	today := Today{
		TimeTotal: getTotalWorkoutTime(todaysWorkout),
	}

	return HevyData{
		Goal:    goal,
		AllTime: allTime,
		Year:    year,
		Week:    week,
		Today:   today,
	}, nil
}

func getYearWorkouts(workouts Workouts) Workouts {
	currentYear, _ := time.Now().ISOWeek()

	var thisYearWorkouts Workouts
	for _, workout := range workouts {
		year := workout.CreatedAt.Year()
		if currentYear == year {
			thisYearWorkouts = append(thisYearWorkouts, workout)
		}
	}

	return thisYearWorkouts
}

func getWeekWorkouts(workouts Workouts) Workouts {
	_, currentWeek := time.Now().ISOWeek()

	var thisYearWorkouts Workouts
	for _, workout := range workouts {
		_, week := workout.CreatedAt.ISOWeek()
		if currentWeek == week {
			thisYearWorkouts = append(thisYearWorkouts, workout)
		}
	}

	return thisYearWorkouts
}

func getTodayWorkout(workouts Workouts) Workouts {
	today := time.Now()
	var todayWorkouts Workouts

	for _, workout := range workouts {
		if workout.CreatedAt.Year() == today.Year() &&
			workout.CreatedAt.Month() == today.Month() &&
			workout.CreatedAt.Day() == today.Day() {
			todayWorkouts = append(todayWorkouts, workout)
		}
	}

	return todayWorkouts
}

func getTotalWorkouts(workouts Workouts) int {
	return len(workouts)
}

func getWorkoutTime(workout Workout) int {
	duration := workout.EndTime.Sub(workout.StartTime)
	return int(duration.Seconds())
}

func getTotalWorkoutTime(workouts Workouts) int {
	totalTime := 0
	for _, workout := range workouts {
		totalTime += getWorkoutTime(workout)
	}
	return totalTime
}

func getGoalMet(goal int, weekWorkouts Workouts) bool {
	return len(weekWorkouts) >= goal
}

func getSumGoalMet(goal int, workouts Workouts) int {
	yearWeekWorkoutsMap := groupWorkoutsByWeek(workouts)
	goalMetCount := 0
	for _, weekWorkouts := range yearWeekWorkoutsMap {
		if getGoalMet(goal, weekWorkouts) {
			goalMetCount++
		}
	}
	return goalMetCount
}

func groupWorkoutsByWeek(workouts Workouts) map[string][]Workout {
	weekWorkoutsMap := make(map[string][]Workout)

	for _, workout := range workouts {
		year, week := workout.CreatedAt.ISOWeek()
		key := fmt.Sprintf("%d-%02d", year, week)
		weekWorkoutsMap[key] = append(weekWorkoutsMap[key], workout)
	}

	return weekWorkoutsMap
}
