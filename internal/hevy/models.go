package hevy

import "time"

type Workout struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Exercises   []Exercise `json:"exercises"`
}

type Exercise struct {
	Index              int    `json:"index"`
	Title              string `json:"title"`
	Notes              string `json:"notes"`
	ExerciseTemplateID string `json:"exercise_template_id"`
	SupersetsID        *int   `json:"supersets_id"`
	Sets               []Set  `json:"sets"`
}

type Set struct {
	Index           int      `json:"index"`
	Type            string   `json:"type"` // normal, warmup, dropset, failure
	WeightKg        *float64 `json:"weight_kg"`
	Reps            *int     `json:"reps"`
	DistanceMeters  *float64 `json:"distance_meters"`
	DurationSeconds *float64 `json:"duration_seconds"`
	RPE             *float64 `json:"rpe"`
}

type Workouts []Workout
type WorkoutResponse struct {
	Page      int      `json:"page"`
	PageCount int      `json:"page_count"`
	Workouts  Workouts `json:"workouts"`
}

type HevyData struct {
	Goal    int     `json:"goal"`
	Today   Today   `json:"today"`
	Week    Week    `json:"week"`
	Year    Year    `json:"year"`
	AllTime AllTime `json:"all_time"`
}

type Today struct {
	TimeTotal int `json:"time_total"`
	// TotalVolume int `json:"total_volume"`
}

type Week struct {
	TotalWorkouts int  `json:"total_workouts"`
	TimeTotal     int  `json:"time_total"`
	GoalMet       bool `json:"goal_met"`
	// TotalVolume   int  `json:"total_volume"`
}

type Year struct {
	TotalWorkouts int `json:"total_workouts"`
	TimeTotal     int `json:"time_total"`
	WeeksGoalMet  int `json:"weeks_goal_met"`
	// LongestGoalStreak int `json:"longest_goal_streak"`
	// LongestStreak int `json:"longest_streak"`
	// TotalVolume   int `json:"total_volume"`
}

type AllTime struct {
	TotalWorkouts int `json:"total_workouts"`
	TimeTotal     int `json:"time_total"`
	WeeksGoalMet  int `json:"weeks_goal_met"`
	// LongestGoalStreak int `json:"longest_goal_streak"`
	// LongestStreak     int `json:"longest_streak"`
	// TotalVolume       int `json:"total_volume"`
}
