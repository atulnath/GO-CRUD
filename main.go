package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS people (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	//create a router
	r := mux.NewRouter()
	r.HandleFunc("/people", getUsers(db)).Methods("GET")
	r.HandleFunc("/people/{id}", getUser(db)).Methods("GET")
	r.HandleFunc("/people", createUser(db)).Methods("POST")
	r.HandleFunc("/people/{id}", updateUser(db)).Methods("PUT")
	r.HandleFunc("/people/{id}", deleteUser(db)).Methods("DELETE")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", jsonContentTypeMiddleware(r)))
}

// getUsers handles the GET request to fetch all users
// from the database and return them as JSON.
func getUsers(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email FROM people")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Error fetching users")
			return
		}
		defer rows.Close()

		var people []Person
		for rows.Next() {
			var p Person
			if err := rows.Scan(&p.ID, &p.Name, &p.Email); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("Error scanning user")
				return
			}
			people = append(people, p)
		}
		json.NewEncoder(w).Encode(people)
	}
}

// getUser handles the GET request to fetch a single user
// from the database by ID and return it as JSON.
// It uses the gorilla/mux package to extract the ID from the URL.
func getUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		row := db.QueryRow("SELECT id, name, email FROM people WHERE id = $1", id)
		var p Person
		if err := row.Scan(&p.ID, &p.Name, &p.Email); err != nil {
			//TODO: handle error
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
		}
		json.NewEncoder(w).Encode(p)
	}
}

func createUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Person
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO people (name, email) VALUES ($1, $2) RETURNING id", p.Name, p.Email).Scan(&p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(p)
	}
}

func updateUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var p Person
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE people SET name = $1, email = $2 WHERE id = $3", p.Name, p.Email, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func deleteUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		result, err := db.Exec("DELETE FROM people WHERE id = $1", id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rowsChanged, err := result.RowsAffected()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if rowsChanged == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode("User deleted")
	}
}
func jsonContentTypeMiddleware(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handle.ServeHTTP(w, r)
	})
}
