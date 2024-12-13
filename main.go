package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
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

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	router.HandleFunc("POST /news", func(w http.ResponseWriter, r *http.Request) {
		output, _ := parallelScraper()
		w.Header().Set("Content-Type", "application/json")
		for s, items := range output {
			err := insertData(items, s)
			if err != nil {
				fmt.Println(err)
			}
		}
		err := json.NewEncoder(w).Encode(output)
		if err != nil {
			fmt.Println(err)
		}
	})

	router.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		output, err := getLastArticles()
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			fmt.Println(err)
		}
	})

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
