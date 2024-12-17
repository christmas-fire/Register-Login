package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/christmas-fire/register-login/internal/users"
	"golang.org/x/crypto/bcrypt"
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

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u users.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		err = users.DeleteUser(db, u.Username)
		if err != nil {
			http.Error(w, "can't delete user: %w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Printf("user: '%s' has deleted", u.Username)
	}
}

func UpdateUserPasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Структура для входных данных
		var req struct {
			Username        string `json:"username"`
			CurrentPassword string `json:"current_password"`
			NewPassword     string `json:"new_password"`
		}

		// Декодирование JSON
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid data format", http.StatusBadRequest)
			return
		}

		// Проверка обязательных полей
		if req.Username == "" || req.CurrentPassword == "" || req.NewPassword == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Получение текущего хэша пароля из базы данных
		var storedPasswordHash string
		query := `SELECT password FROM users WHERE username = $1`
		err = db.QueryRow(query, req.Username).Scan(&storedPasswordHash)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Проверка текущего пароля
		err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(req.CurrentPassword))
		if err != nil {
			http.Error(w, "Invalid current password", http.StatusUnauthorized)
			return
		}

		// Хэширование нового пароля
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
			return
		}

		// Обновление пароля в базе данных
		updateQuery := `UPDATE users SET password = $1 WHERE username = $2`
		_, err = db.Exec(updateQuery, newPasswordHash, req.Username)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusNoContent)
		log.Printf("User's password updated: username='%s'", req.Username)
	}
}

func UpdateUserUsernameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username    string `json:"username"`
			NewUsername string `json:"new_username"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid data format", http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.NewUsername == "" {
			http.Error(w, "missing required fields: 'username' or 'new_username'", http.StatusBadRequest)
			return
		}

		err = users.UpdateUserUsername(db, req.Username, req.NewUsername)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update username: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Printf("user's username updated: old='%s', new='%s'", req.Username, req.NewUsername)
	}
}
