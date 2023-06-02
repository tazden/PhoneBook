package repo

import (
	"database/sql"
	"fmt"
	"github.com/DenisTaztdinov/PhoneBook/internal/entity"
)

type ContactsRepoImpl struct {
	Db *sql.DB
}

func NewContactsRepo(db *sql.DB) *ContactsRepoImpl {
	return &ContactsRepoImpl{
		Db: db,
	}
}

func (r *ContactsRepoImpl) GetAll() ([]entity.Contact, error) {
	rows, err := r.Db.Query("SELECT id, first_name, last_name, phone, email FROM contacts")
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
func (r *ContactsRepoImpl) GetByID(id int) (*entity.Contact, error) {
	row := r.Db.QueryRow("SELECT id, first_name, last_name, phone, email FROM contacts WHERE id = $1", id)

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
func (r *ContactsRepoImpl) Create(contact *entity.Contact) error {
	_, err := r.Db.Exec("INSERT INTO contacts (first_name, last_name, phone, email) VALUES ($1, $2, $3, $4)",
		contact.FirstName, contact.LastName, contact.Phone, contact.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *ContactsRepoImpl) Update(contact *entity.Contact) error {
	_, err := r.Db.Exec("UPDATE contacts SET first_name = $1, last_name = $2, phone = $3, email = $4 WHERE id = $5",
		contact.FirstName, contact.LastName, contact.Phone, contact.Email, contact.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет контакт по указанному идентификатору
func (r *ContactsRepoImpl) Delete(id int) error {
	_, err := r.Db.Exec("DELETE FROM contacts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
