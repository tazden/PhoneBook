package internal

import "github.com/DenisTaztdinov/PhoneBook/internal/entity"

type ContactRepository interface {
	GetAll() ([]entity.Contact, error)
	GetByID(id int) (*entity.Contact, error)
	Create(contact *entity.Contact) error
	Update(contact *entity.Contact) error
	Delete(id int) error
}
