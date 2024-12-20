package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/christmas-fire/register-login/internal/models"
	"github.com/christmas-fire/register-login/internal/repository/postgres"
)

func GetAllTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := postgres.GetAllUsers(db)
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
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		err = postgres.AddUser(db, u.Username, u.Email, u.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
		log.Printf("user: '%s' has registered", u.Username)
	}
}

func ValidateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		err = postgres.ValidateUser(db, u.Username, u.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Printf("user: '%s' has logined", u.Username)
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		err = postgres.DeleteUser(db, u.Username)
		if err != nil {
			http.Error(w, "can't delete user: %w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Printf("user: '%s' has deleted", u.Username)
	}
}
