package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/henryppercy/accountability-api/internal/db"
	"github.com/henryppercy/accountability-api/internal/ds"
	"github.com/henryppercy/accountability-api/internal/hevy"
)

func main() {
	dbs, err := db.InitDatabase("./database/accountability.sqlite")

	if err != nil {
		log.Fatal("error initializing DB connection: ", err)
	}

	err = dbs.Ping()
	if err != nil {
		log.Fatal("error initializing DB connection: ping error: ", err)
	}

	fmt.Println("database initialized..")

	initialized, err := db.IsDatabaseInitialized()
	if err != nil {
		log.Fatal("error checking database initialization: ", err)
	}

	if initialized {
		fmt.Println("Database already initialized. Skipping data insertion.")
	} else {
		workouts, err := hevy.FetchWorkouts()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, workout := range workouts {
			err = hevy.InsertWorkout(dbs, &workout)
			if err != nil {
				log.Fatal("error inserting test data: ", err)
			}
			fmt.Println("workout inserted", workout.Title)
		}

		userData, err := ds.FetchUserData()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = ds.InsertUserData(dbs, &userData)
		if err != nil {
			log.Fatal("error inserting test data: ", err)
		}
		fmt.Println("DS data inserted")

		dailyTimes, err := ds.FetchDailyTimes()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	
		for _, dailyTime := range dailyTimes {
			err = ds.InsertDailyTime(dbs, &dailyTime)
			if err != nil {
				log.Fatal("error inserting test data: ", err)
			}
			fmt.Println("DS daily time inserted", dailyTime.Date)
		}

		err = db.MarkDatabaseAsInitialized()
		if err != nil {
			log.Fatal("error marking database as initialized: ", err)
		}
	}

	DSData, err := ds.GetDSData(dbs)

	if err != nil {
		log.Fatal("error retrieving DS data from db: ", err)
	}

	jsonData, _ := json.MarshalIndent(DSData, "", "  ")
	fmt.Println(string(jsonData))

	hevyData, err := hevy.GetHevyData(dbs)

	if err != nil {
		log.Fatal("error retrieving DS data from db: ", err)
	}

	jsonData, _ = json.MarshalIndent(hevyData, "", "  ")
	fmt.Println(string(jsonData))


	pageData := PageData{
		DSData:   DSData,
		HevyData: hevyData,
		DaysOfWeek: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderHomePage(w, pageData)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type PageData struct {
	DSData    ds.DSData
	HevyData  hevy.HevyData
	DaysOfWeek   []string
}

func renderHomePage(w http.ResponseWriter, pageData PageData) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
