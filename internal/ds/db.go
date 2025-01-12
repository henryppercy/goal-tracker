package ds

import (
	"database/sql"
	"fmt"
)

func InsertDailyTime(db *sql.DB, dailyTime *DailyTime) error {
    stmt, err := db.Prepare("INSERT INTO ds_daily_times(date, user_id, time_seconds, goal_reached) VALUES(?, ?, ?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare daily time insert: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(dailyTime.Date, dailyTime.UserID, dailyTime.TimeSeconds, dailyTime.GoalReached)
    if err != nil {
        return fmt.Errorf("failed to insert daily time: %w", err)
    }

    return nil
}

func InsertUserData(db *sql.DB, userData *UserData) error {
	stmt, err := db.Prepare(`INSERT INTO ds_user_data(
        watch_time,
        daily_goal_seconds
    ) VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare UserData insert: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		userData.WatchTime,
		userData.DailyGoalSeconds,
	)
	if err != nil {
		return fmt.Errorf("failed to insert UserData: %w", err)
	}

	return nil
}

func GetDailyTimes(db *sql.DB) (DailyTimes, error) {
	rows, err := db.Query("SELECT date, user_id, time_seconds, goal_reached FROM ds_daily_times")
	if err != nil {
		return nil, fmt.Errorf("failed to query daily times: %w", err)
	}
	defer rows.Close()

	var dailyTimes DailyTimes
	for rows.Next() {
		var dailyTime DailyTime
		err := rows.Scan(&dailyTime.Date, &dailyTime.UserID, &dailyTime.TimeSeconds, &dailyTime.GoalReached)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily time row: %w", err)
		}
		dailyTimes = append(dailyTimes, dailyTime)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return dailyTimes, nil
}

func GetUserData(db *sql.DB) (*UserData, error) {
	row := db.QueryRow("SELECT watch_time, daily_goal_seconds FROM ds_user_data")

	var userData UserData
	err := row.Scan(&userData.WatchTime, &userData.DailyGoalSeconds)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user data found")
		}
		return nil, fmt.Errorf("failed to scan user data row: %w", err)
	}

	return &userData, nil
}
