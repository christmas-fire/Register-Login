package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/christmas-fire/register-login/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Проверяем входные данные от пользователя
func ValidateUserData(username, email, password string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must have at least 3 characters")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must have at least 8 characters")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// Регистрация нового пользователя в базе данных
func AddUser(db *sql.DB, username, email, password string) error {
	if err := ValidateUserData(username, email, password); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}

	checkQuery := `
        SELECT EXISTS (
            SELECT 1
            FROM users
            WHERE username = $1 OR email = $2
        )
    `
	var alreadyExists bool
	if err := db.QueryRow(checkQuery, username, email).Scan(&alreadyExists); err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}

	if alreadyExists {
		return fmt.Errorf("user with the same username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	insertQuery := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
    `
	if _, err := db.Exec(insertQuery, username, email, hashedPassword); err != nil {
		return fmt.Errorf("error inserting new user into the database: %w", err)
	}

	return nil
}

// Авторизация пользователя
func ValidateUser(db *sql.DB, username, password string) error {
	var hashedPassword string

	query := `
        SELECT password FROM users WHERE username = $1
    `
	err := db.QueryRow(query, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("error retrieving user from the database: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

// Удаление пользователя
func DeleteUser(db *sql.DB, username string) error {
	query := `
		DELETE FROM users
		WHERE username = $1
	`

	_, err := db.Exec(query, username)
	if err != nil {
		return fmt.Errorf("can't delete user '%s': %v", username, err)
	}

	return nil
}

// Получить всех пользователей
func GetAllUsers(db *sql.DB) ([]models.User, error) {
	query := `
		SELECT username, email, password FROM users
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Username, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
