package handlers

import (
	"TodoApp/models"
	"encoding/json"
	"net/http"
)

func GetTODOList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := models.DB.Query("SELECT id, title, status, priority FROM todos")
	if err != nil {
		http.Error(w, "database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []models.TodoModel
	for rows.Next() {
		var item models.TodoModel
		if err := rows.Scan(&item.ID, &item.Title, &item.Status, &item.Priority); err != nil {
			http.Error(w, "error scanning rows", http.StatusInternalServerError)
			return
		}
		list = append(list, item)
	}

	if list == nil {
		list = []models.TodoModel{}
	}

	json.NewEncoder(w).Encode(list)
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	var item models.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "failed to parse data", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO todos (title, status, priority) VALUES ($1, $2, $3) RETURNING id"
	err := models.DB.QueryRow(query, item.Title, item.Status, item.Priority).Scan(&item.ID)
	if err != nil {
		http.Error(w, "failed to save data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM todos WHERE id = $1"
	result, err := models.DB.Exec(query, idParam)
	if err != nil {
		http.Error(w, "failed to delete item", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "item deleted successfully"})
}

type StatusUpdateRequest struct {
	Status string `json:"status"`
}

func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "missing id param", http.StatusBadRequest)
		return
	}

	var request StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	query := "UPDATE todos SET status = $1 WHERE id = $2"
	result, err := models.DB.Exec(query, request.Status, idParam)
	if err != nil {
		http.Error(w, "failed to update database", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "status update success"})
}

func HandleTodo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTODOList(w, r)
	case http.MethodPost:
		AddTodo(w, r)
	case http.MethodDelete:
		DeleteTodo(w, r)
	case http.MethodPut:
		UpdateStatus(w, r)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}
