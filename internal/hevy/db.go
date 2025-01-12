package hevy

import (
	"database/sql"
	"fmt"
	"time"
)

func InsertWorkout(db *sql.DB, workout *Workout) error {
    stmt, err := db.Prepare("INSERT INTO workouts(id, title, description, start_time, end_time, updated_at, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare workout insert: %w", err)
    }
    _, err = stmt.Exec(workout.ID, workout.Title, workout.Description, workout.StartTime, workout.EndTime, workout.UpdatedAt, workout.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to insert workout: %w", err)
    }

    exerciseStmt, err := db.Prepare("INSERT INTO exercises(workout_id, \"index\", title, notes, exercise_template_id, supersets_id) VALUES(?, ?, ?, ?, ?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare exercise insert: %w", err)
    }
    defer exerciseStmt.Close()

    setStmt, err := db.Prepare("INSERT INTO sets(exercise_id, \"index\", type, weight_kg, reps, distance_meters, duration_seconds, rpe) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare set insert: %w", err)
    }
    defer setStmt.Close()

    for _, exercise := range workout.Exercises {
        result, err := exerciseStmt.Exec(workout.ID, exercise.Index, exercise.Title, exercise.Notes, exercise.ExerciseTemplateID, exercise.SupersetsID)
        if err != nil {
            return fmt.Errorf("failed to insert exercise: %w", err)
        }

        exerciseID, err := result.LastInsertId()
        if err != nil {
            return fmt.Errorf("failed to retrieve exercise ID: %w", err)
        }

        for _, set := range exercise.Sets {
            _, err := setStmt.Exec(exerciseID, set.Index, set.Type, set.WeightKg, set.Reps, set.DistanceMeters, set.DurationSeconds, set.RPE)
            if err != nil {
                return fmt.Errorf("failed to insert set: %w", err)
            }
        }
    }

    return nil
}

func GetWorkouts(db *sql.DB) (Workouts, error) {
	workoutRows, err := db.Query("SELECT id, title, description, start_time, end_time, updated_at, created_at FROM workouts")
	if err != nil {
		return nil, fmt.Errorf("failed to query workouts: %w", err)
	}
	defer workoutRows.Close()

	var workouts []Workout

	for workoutRows.Next() {
		var workout Workout
		var startTime, endTime, updatedAt, createdAt string

		err := workoutRows.Scan(&workout.ID, &workout.Title, &workout.Description, &startTime, &endTime, &updatedAt, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout row: %w", err)
		}

        workout.StartTime, err = time.Parse("2006-01-02 15:04:05 -0700 MST", startTime)
        if err != nil {
            return nil, fmt.Errorf("failed to parse workout start_time: %w", err)
        }
        workout.EndTime, err = time.Parse("2006-01-02 15:04:05 -0700 MST", endTime)
        if err != nil {
            return nil, fmt.Errorf("failed to parse workout end_time: %w", err)
        }
        workout.UpdatedAt, err = time.Parse("2006-01-02 15:04:05 -0700 MST", updatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to parse workout updated_at: %w", err)
        }
        workout.CreatedAt, err = time.Parse("2006-01-02 15:04:05 -0700 MST", createdAt)
        if err != nil {
            return nil, fmt.Errorf("failed to parse workout created_at: %w", err)
        }        

		workout.Exercises, err = getExercisesForWorkout(db, workout.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch exercises for workout %s: %w", workout.ID, err)
		}

		workouts = append(workouts, workout)
	}

	return workouts, nil
}

func getExercisesForWorkout(db *sql.DB, workoutID string) ([]Exercise, error) {
	exerciseRows, err := db.Query("SELECT id, \"index\", title, notes, exercise_template_id, supersets_id FROM exercises WHERE workout_id = ?", workoutID)
	if err != nil {
		return nil, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer exerciseRows.Close()

	var exercises []Exercise

	for exerciseRows.Next() {
		var exercise Exercise
		var exerciseID int

		err := exerciseRows.Scan(&exerciseID, &exercise.Index, &exercise.Title, &exercise.Notes, &exercise.ExerciseTemplateID, &exercise.SupersetsID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise row: %w", err)
		}

		exercise.Sets, err = getSetsForExercise(db, exerciseID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch sets for exercise %d: %w", exerciseID, err)
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func getSetsForExercise(db *sql.DB, exerciseID int) ([]Set, error) {
	setRows, err := db.Query("SELECT \"index\", type, weight_kg, reps, distance_meters, duration_seconds, rpe FROM sets WHERE exercise_id = ?", exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sets: %w", err)
	}
	defer setRows.Close()

	var sets []Set

	for setRows.Next() {
		var set Set

		err := setRows.Scan(&set.Index, &set.Type, &set.WeightKg, &set.Reps, &set.DistanceMeters, &set.DurationSeconds, &set.RPE)
		if err != nil {
			return nil, fmt.Errorf("failed to scan set row: %w", err)
		}

		sets = append(sets, set)
	}

	return sets, nil
}
