package users

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AddUser(db *sql.DB, username, email, password string) error {
	checkQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM USERS
			WHERE username = $1 OR email = $2
		)
	`

	var alreadyExists bool
	err := db.QueryRow(checkQuery, username, email).Scan(&alreadyExists)
	if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if alreadyExists {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	insertQuery := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	`

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = db.Exec(insertQuery, username, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("error insert new user: %w", err)
	}

	return nil
}

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
		return fmt.Errorf("error validate user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func UpdateUserPassword(db *sql.DB, username, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	// SQL-запрос на обновление
	query := `
		UPDATE users
		SET password = $2
		WHERE username = $1
	`

	// Выполнение запроса
	_, err = db.Exec(query, username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("can't update user's '%s' password: %v", username, err)
	}

	return nil
}

// UpdateUserUsername обновляет имя пользователя
func UpdateUserUsername(db *sql.DB, currentUsername, newUsername string) error {
	// SQL-запрос на обновление
	query := `
		UPDATE users
		SET username = $2
		WHERE username = $1
	`

	// Выполнение запроса
	_, err := db.Exec(query, currentUsername, newUsername)
	if err != nil {
		return fmt.Errorf("can't update username from '%s' to '%s': %v", currentUsername, newUsername, err)
	}

	return nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
		SELECT username, email, password FROM users
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
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
