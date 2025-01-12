package ds

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/henryppercy/accountability-api/internal/format"
)

func GetDSData(db *sql.DB) (DSData, error) {
	userData, err := GetUserData(db)
	if err != nil {
		return DSData{}, fmt.Errorf("failed to fetch user data: %w", err)
	}

	dailyTimes, err := GetDailyTimes(db)
	if err != nil {
		return DSData{}, fmt.Errorf("failed to fetch daily times: %w", err)
	}

	var todayTimeSeconds format.Seconds = 0
	todayGoalMet := false

	for _, day := range dailyTimes {
		if day.Date == time.Now().Format("2006-01-02") {
			todayTimeSeconds = day.TimeSeconds
			todayGoalMet = day.GoalReached
			break
		}
	}

	platformStreak := calculateStreak(dailyTimes, func(day DailyTime) bool {
		return day.TimeSeconds > 0 // streak using platform
	})
	goalStreak := calculateStreak(dailyTimes, func(day DailyTime) bool {
		return day.GoalReached // streak completing goal
	})

	daysThisWeek := amountOfDaysGoalHitThisWeek(dailyTimes)
	weekGoalStatus := goalStatusPerDay(dailyTimes, userData.DailyGoalSeconds)

	return DSData{
		TotalWatchTime:   userData.WatchTime,
		DailyGoal: userData.DailyGoalSeconds,
		TodayTime: todayTimeSeconds,
		WeekTime: timeTotalThisWeek(dailyTimes),
		TodayGoalMet:     todayGoalMet,
		PlatformStreak:   platformStreak,
		GoalStreak:       goalStreak,
		GoalsTotalWeek:   daysThisWeek,
		WeekGoalStatus:   weekGoalStatus,
	}, nil
}

// calculates the current streak for a given condition
func calculateStreak(dailyTimes DailyTimes, condition func(DailyTime) bool) int {
	streak := 0
	sort.Slice(dailyTimes, func(i, j int) bool {
		return dailyTimes[i].Date > dailyTimes[j].Date // Sort descending by date
	})
	for _, day := range dailyTimes {
		if condition(day) {
			streak++
		} else {
			break
		}
	}
	return streak
}

func timeTotalThisWeek(dailyTimes DailyTimes) format.Seconds {
	currentYear, currentWeek := time.Now().ISOWeek()
	var timeThisWeek format.Seconds = 0

	for _, day := range dailyTimes {
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			continue
		}
		year, week := date.ISOWeek()
		if year == currentYear && week == currentWeek {
			timeThisWeek += day.TimeSeconds
		}
	}
	return timeThisWeek
}

// counts how many days in the current week have been completed
func amountOfDaysGoalHitThisWeek(dailyTimes DailyTimes) int {
	currentYear, currentWeek := time.Now().ISOWeek()
	daysThisWeek := 0

	for _, day := range dailyTimes {
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			continue
		}
		year, week := date.ISOWeek()
		if year == currentYear && week == currentWeek && day.GoalReached {
			daysThisWeek++
		}
	}
	return daysThisWeek
}


func goalStatusPerDay(dailyTimes DailyTimes, dailyGoalSeconds format.Seconds) WeekGoalStatus {
    weekGoalStatus := WeekGoalStatus{
        "Mon": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
        "Tue": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
        "Wed": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
        "Thu": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
        "Fri": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
        "Sat": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
        "Sun": {GoalMet: false, PercentOfGoal: 0, PercentOfGoalLimit: 0},
    }

    currentYear, currentWeek := time.Now().ISOWeek()

    for _, day := range dailyTimes {
        date, err := time.Parse("2006-01-02", day.Date)
        if err != nil {
            log.Printf("Error parsing date %s: %v", day.Date, err)
            continue
        }

        year, week := date.ISOWeek()
        if year == currentYear && week == currentWeek {
            dayOfWeek := date.Weekday().String()[:3]

			percentOfGoalLimit := int(math.Min(100, math.Round(float64(day.TimeSeconds)/float64(dailyGoalSeconds)*100)))
			percentOfGoal := int(math.Round(float64(day.TimeSeconds) / float64(dailyGoalSeconds) * 100))

            weekGoalStatus[dayOfWeek] = DayStats{
                GoalMet:       day.GoalReached,
                PercentOfGoal: percentOfGoal,
				PercentOfGoalLimit: percentOfGoalLimit,
            }
        }
    }

    return weekGoalStatus
}
