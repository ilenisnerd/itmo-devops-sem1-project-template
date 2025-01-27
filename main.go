package main

import (
	"log"
	"net/http"
	_ "os"
	"project-sem-1/database"
	"project-sem-1/handlers"

	"github.com/gorilla/mux"
)

func main() {

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/v0/prices", handlers.PostPrices(db)).Methods("POST")
	r.HandleFunc("/api/v0/prices", handlers.GetPrices(db)).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
