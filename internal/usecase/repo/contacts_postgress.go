package repo

import (
	"database/sql"
	"fmt"
	"github.com/DenisTaztdinov/PhoneBook/internal/entity"
)

type PostgresSQLRepository struct {
	db *sql.DB
}

func (r *PostgresSQLRepository) GetAll() ([]entity.Contact, error) {
	rows, err := r.db.Query("SELECT id, first_name, last_name, phone, email FROM contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []entity.Contact{}
	for rows.Next() {
		var contact entity.Contact
		err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Phone, &contact.Email)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetByID возвращает контакт по указанному идентификатору
func (r *PostgresSQLRepository) GetByID(id int) (*entity.Contact, error) {
	row := r.db.QueryRow("SELECT id, first_name, last_name, phone, email FROM contacts WHERE id = $1", id)

	var contact entity.Contact
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
func (r *PostgresSQLRepository) Create(contact *entity.Contact) error {
	_, err := r.db.Exec("INSERT INTO contacts (first_name, last_name, phone, email) VALUES ($1, $2, $3, $4)",
		contact.FirstName, contact.LastName, contact.Phone, contact.Email)
	if err != nil {
		return err
	}

	return nil
}
