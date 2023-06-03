package app

import (
	"database/sql"
	"fmt"
	"github.com/DenisTaztdinov/PhoneBook/config"
	"github.com/DenisTaztdinov/PhoneBook/internal/usecase"
	"github.com/DenisTaztdinov/PhoneBook/internal/usecase/repo"
	"github.com/DenisTaztdinov/PhoneBook/pkg/logger"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	db, err := sql.Open("postgres", cfg.PG.URL)
	if err != nil {
		l.Fatal(err)
	}
	defer db.Close()
	repo := repo.NewContactsRepo(db)
	handler := &usecase.ContactHandler{Repo: repo}

	http.HandleFunc("/contacts", handler.GetAllContacts)

	port := cfg.HTTP.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	l.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
