package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase(dbPath string) (*sql.DB ,error) {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	err = createInitStatusTable()
	if err != nil {
		return nil, fmt.Errorf("error creating init status table: %v", err)
	}

	err = createHevyTables()
	if err != nil {
		return nil, fmt.Errorf("error creating Hevy tables: %v", err)
	}

    err = createDSTables()
	if err != nil {
		return nil, fmt.Errorf("error creating DS tables: %v", err)
	}
	return db, nil
}

func createHevyTables() error {
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS workouts (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT,
            start_time TEXT NOT NULL,
            end_time TEXT NOT NULL,
            updated_at TEXT NOT NULL,
            created_at TEXT NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS exercises (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            workout_id TEXT NOT NULL,
            "index" INTEGER NOT NULL,
            title TEXT NOT NULL,
            notes TEXT,
            exercise_template_id TEXT NOT NULL,
            supersets_id INTEGER,
            FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS sets (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            exercise_id INTEGER NOT NULL,
            "index" INTEGER NOT NULL,
            type TEXT NOT NULL,
            weight_kg REAL,
            reps INTEGER,
            distance_meters REAL,
            duration_seconds REAL,
            rpe REAL,
            FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
        );`,
	}

	for _, query := range createTableQueries {
		_, err := db.ExecContext(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDSTables() error {
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS ds_daily_times (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            user_id TEXT NOT NULL,
            time_seconds INTEGER NOT NULL,
            goal_reached BOOLEAN NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS ds_user_data (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            watch_time INTEGER NOT NULL,
            daily_goal_seconds INTEGER NOT NULL
        );`,
	}

	for _, query := range createTableQueries {
		_, err := db.ExecContext(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}

func createInitStatusTable() error {
	query := `CREATE TABLE IF NOT EXISTS init_status (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		is_initialized BOOLEAN NOT NULL
	);`

	_, err := db.ExecContext(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error creating init_status table: %v", err)
	}

	row := db.QueryRow("SELECT is_initialized FROM init_status LIMIT 1")

	var isInitialized bool
	err = row.Scan(&isInitialized)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := db.Exec("INSERT INTO init_status (is_initialized) VALUES (false)")
			if err != nil {
				return fmt.Errorf("error initializing init_status table: %v", err)
			}
		} else {
			return fmt.Errorf("error checking init_status: %v", err)
		}
	}

	return nil
}

func IsDatabaseInitialized() (bool, error) {
	var isInitialized bool
	row := db.QueryRow("SELECT is_initialized FROM init_status LIMIT 1")
	err := row.Scan(&isInitialized)
	if err != nil {
		return false, fmt.Errorf("error checking database initialization: %v", err)
	}
	return isInitialized, nil
}

func MarkDatabaseAsInitialized() error {
	_, err := db.Exec("UPDATE init_status SET is_initialized = true")
	if err != nil {
		return fmt.Errorf("error marking database as initialized: %v", err)
	}
	return nil
}
