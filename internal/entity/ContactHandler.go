package entity

import (
	"encoding/json"
	"log"
	"net/http"
)

type ContactHandler struct {
	Repo ContactRepository
}

func (h *ContactHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.Repo.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(contacts)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
