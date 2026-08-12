package main

import ( // FELT LIKE COMMITMAXXING
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

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError) // ping http for a scan error.
			return 
			}
		 items = append(items, item)
		}
	json.NewEncoder(w).Encode(items)
}

// GET /items/{id} Fetch a item by ID
func handleGetItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id") // Extract {id} from str

	var item Item
	err := dbPool.QueryRow(r.Context(), "SELECT id, name, price FROM items WHERE id = $1", id).
		Scan(&item.ID, &item.Name, &item.Price)
	// adding error if else statements
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
} // finishing the handleGetItembyID function off.

// POST /items Add a new item to inventory very easy access
func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //  handleCreateItem

	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest) // http error handling
		return
	}

	_, err := dbPool.Exec(r.Context(), "INSERT INTO items (id, name, price) VALUES ($1, $2, $3)", // use sql to insert items.
		newItem.ID, newItem.Name, newItem.Price)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError) // error handling if sql doesnt work.
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}
// end of handleCreateItem function

// PUT /items/{id} update an existing item
func handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // reroute
	id := r.PathValue("id")

	var updatedFields Item // begin updatedFields variable
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest) // http error handling
		return
	}
