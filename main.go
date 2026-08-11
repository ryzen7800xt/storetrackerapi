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
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Supabase connection failed.")
	}
	var err error
	dbpool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()
}
