package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

var db *sql.DB

func main() {
	router := http.NewServeMux()

	var err error
	connStr := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)

	go scrapingScheduler()

	router.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		output, err := getLastArticles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}

func scrapingScheduler() {
	runScraper()

	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runScraper()
		}
	}
}

func runScraper() {
	output, err := parallelScraper()
	if err != nil {
		fmt.Printf("Scraping error: %v\n", err)
		return
	}

	for s, items := range output {
		if err := insertData(items, s); err != nil {
			fmt.Printf("Insert error for source %s: %v\n", s, err)
		}
	}
}
