package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/christmas-fire/register-login/internal/users"
)

func GetAllTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := users.GetAllUsers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func AddUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u users.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		err = users.AddUser(db, u.Username, u.Email, u.Password)
		if err != nil {
			http.Error(w, "error add new user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
		log.Printf("user: '%s' has registered", u.Username)
	}
}

func ValidateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u users.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		err = users.ValidateUser(db, u.Username, u.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Printf("user: '%s' has validated", u.Username)
	}
}
