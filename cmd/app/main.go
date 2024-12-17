package main

import (
	"log"
	"net/http"

	"github.com/christmas-fire/register-login/internal/database"
	"github.com/christmas-fire/register-login/internal/rest"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	err := database.CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		rest.AddUserHandler(db)(w, r)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		rest.ValidateUserHandler(db)(w, r)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		rest.GetAllTasksHandler(db)(w, r)
	})

	log.Println("server is running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("error start server:", err)
		return
	}
}
