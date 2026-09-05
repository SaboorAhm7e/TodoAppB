package models

import "database/sql"

type TodoModel struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

var DB *sql.DB
