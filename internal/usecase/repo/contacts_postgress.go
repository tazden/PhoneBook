package repo

import (
	"database/sql"
	"fmt"
)

type PostgresSQLRepository struct {
	db *sql.DB
}

func (r *PostgresSQLRepository) GetAll() ([]Contact, error) {
	rows, err := r.db.Query("SELECT id, first_name, last_name, phone, email FROM contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []Contact{}
	for rows.Next() {
		var contact Contact
		err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Phone, &contact.Email)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetByID возвращает контакт по указанному идентификатору
func (r *PostgresSQLRepository) GetByID(id int) (*Contact, error) {
	row := r.db.QueryRow("SELECT id, first_name, last_name, phone, email FROM contacts WHERE id = $1", id)

	var contact Contact
	err := row.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Phone, &contact.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contact not found")
		}
		return nil, err
	}

	return &contact, nil
}

// Create создает новый контакт
func (r *PostgresSQLRepository) Create(contact *Contact) error {
	_, err := r.db.Exec("INSERT INTO contacts (first_name, last_name, phone, email) VALUES ($1, $2, $3, $4)",
		contact.FirstName, contact.LastName, contact.Phone, contact.Email)
	if err != nil {
		return err
	}

	return nil
}
