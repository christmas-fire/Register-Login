package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/christmas-fire/register-login/internal/config"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	cfg, err := config.LoadConfig("./internal/config/")
	if err != nil {
		log.Fatal(err)
	}

	con := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port, cfg.Sslmode,
	)

	db, err := sql.Open("postgres", con)
	if err != nil {
		log.Fatalf("error connect DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("DB unvaluable: %v", err)
	}

	log.Println("success connect DB")

	return db
}

func CreateTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			jwt TEXT
		)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error create table: %w", err)
	}

	log.Println("table created successfully")
	return nil
}
