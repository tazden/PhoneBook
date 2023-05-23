package entity

type ContactRepository interface {
	GetAll() ([]Contact, error)
	GetByID(id int) (*Contact, error)
	Create(contact *Contact) error
	Update(contact *Contact) error
	Delete(id int) error
}
