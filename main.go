package main

import (
	"TodoApp/handlers"
	"TodoApp/models"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

func main() {
	fmt.Println(" ------ Todo App (Supabase Postgres) --------- ")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		panic("DATABASE_URL environment variable is not set")
	}

	var err error
	models.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer models.DB.Close()

	if err = models.DB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Supabase: %v", err))
	}
	fmt.Println("Connected to Supabase PostgreSQL successfully!")

	// Create table using PostgreSQL DDL syntax (SERIAL instead of AUTOINCREMENT)
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		priority TEXT NOT NULL
	);`

	_, err = models.DB.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/todo", handlers.HandleTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
