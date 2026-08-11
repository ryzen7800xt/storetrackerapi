package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Item struct {
		ID    string  `json:"id"` // Use encoding/json to make the 3 tiers of db sorting
		Name  string  `json:"name"` // name
		Price float64 `json:"price"` // price
}

	var dbpool  *pgxpool.Pool

func main() {
	// first step is to load the supabase connection string
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Supabase connection failed.")
	}
	// second step is to establish za connection pool
	var err error // declare error as a variable
	dbpool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()
	// third step verifing the connection to the db
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v\n", err) // db ping fail log
	}
	fmt.Println(" Successfully connected to Supabase") // verifying connection to the backend.
	// http handling with sql
		http.HandleFunc("GET /items", handleGetItems)
	http.HandleFunc("GET /items/{id}", handleGetItemByID)
	http.HandleFunc("POST /items", handleCreateItem)
	http.HandleFunc("PUT /items/{id}", handleUpdateItem)
	http.HandleFunc("DELETE /items/{id}", handleDeleteItem)

	fmt.Println("Server running on http://localhost:8080...") // log http://localhost location
	log.Fatal(http.ListenAndServe(":8080", nil)) // log a fatal error.
}
// API handling function


// get	/items - fetch inventory from db
func handleGetItems(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json") // set the header file to  app/json

	rows, err := dbPool.Query(r.Context(), "SELECT id, name, price FROM items") // SQL QUERY BACK TO THE GOOFY DB
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError) // me when error
		return
	}
	defer rows.Close()
}
