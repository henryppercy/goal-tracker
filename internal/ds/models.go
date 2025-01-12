package ds

import "github.com/henryppercy/accountability-api/internal/format"

type UserData struct {
	WatchTime        format.Seconds `json:"watchTime"`
	DailyGoalSeconds format.Seconds `json:"dailyGoalSeconds"`
}

type UserDataResponse struct {
	User UserData `json:"user"`
}

type DailyTime struct {
	Date        string         `json:"date"`
	UserID      string         `json:"userId"`
	TimeSeconds format.Seconds `json:"timeSeconds"`
	GoalReached bool           `json:"goalReached"`
}
type DailyTimes []DailyTime

type DSData struct {
	TotalWatchTime format.Seconds `json:"total_watch_time"`
	DailyGoal      format.Seconds `json:"daily_goal"`
	TodayTime      format.Seconds `json:"today_time"`
	WeekTime      format.Seconds `json:"week_time"`
	TodayGoalMet          bool           `json:"today_goal_met"`
	PlatformStreak        int            `json:"platform_streak"`
	GoalStreak            int            `json:"goal_streak"`
	GoalsTotalWeek        int            `json:"goals_total_week"`
	WeekGoalStatus        WeekGoalStatus `json:"week_goal_status"`
}

type WeekGoalStatus map[string]DayStats
type DayStats struct {
	GoalMet            bool `json:"goal_met"`
	PercentOfGoal      int  `json:"percent_of_goal"`
	PercentOfGoalLimit int  `json:"percent_of_goal_limit"`
}
