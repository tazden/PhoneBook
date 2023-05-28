package app

import (
	"database/sql"
	"fmt"
	"github.com/DenisTaztdinov/PhoneBook/config"
	"github.com/DenisTaztdinov/PhoneBook/internal/entity"
	"github.com/DenisTaztdinov/PhoneBook/internal/usecase/repo"
	"github.com/DenisTaztdinov/PhoneBook/pkg/logger"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		l.Fatal(err)
	}
	defer db.Close()

	repo := &repo.PostgresSQLRepository{Db: db}
	handler := &entity.ContactHandler{Repo: repo}

	http.HandleFunc("/contacts", handler.GetAllContacts)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	l.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
