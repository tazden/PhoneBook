package main

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

var validate *validator.Validate

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env.example file")
	}

	router := gin.Default()
	router.GET("/contacts", getContacts)
	router.POST("/contacts", addContact)
	router.GET("/contacts/:id", getContact)
	router.PUT("/contacts/:id", updateContact)
	router.DELETE("/contacts/:id", deleteContact)

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getContacts(c *gin.Context) {
	conn, err := getConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if conn != nil {
			if err := conn.Close(context.Background()); err != nil {
				log.Println("Error closing the database connection:", err)
			}
		}
	}()

	rows, err := conn.Query(context.TODO(), "SELECT * FROM contacts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	contacts := []Contact{}
	for rows.Next() {
		contact := Contact{}
		err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		contacts = append(contacts, contact)
	}

	c.JSON(http.StatusOK, contacts)
}

func addContact(c *gin.Context) {
	var contact Contact
	validate = validator.New()
	if err := validate.Struct(contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := getConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if conn != nil {
			if err := conn.Close(context.Background()); err != nil {
				log.Println("Error closing the database connection:", err)
			}
		}
	}()

	err = conn.QueryRow(context.TODO(),
		"INSERT INTO contacts (first_name, last_name, phone_number, email) VALUES ($1, $2, $3, $4) RETURNING id",
		contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Email).Scan(&contact.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}
func getContact(c *gin.Context) {
	id := c.Param("id")

	conn, err := getConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if conn != nil {
			if err := conn.Close(context.Background()); err != nil {
				log.Println("Error closing the database connection:", err)
			}
		}
	}()

	contact := Contact{}
	err = conn.QueryRow(context.TODO(),
		"SELECT * FROM contacts WHERE id = $1",
		id).Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}
func updateContact(c *gin.Context) {
	id := c.Param("id")

	var contact Contact
	validate = validator.New()
	if err := validate.Struct(contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := getConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if conn != nil {
			if err := conn.Close(context.Background()); err != nil {
				// Handle the error appropriately, such as logging or returning an error message.
				log.Println("Error closing the database connection:", err)
			}
		}
	}()

	_, err = conn.Exec(context.TODO(),
		"UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, email = $4 WHERE id = $5",
		contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}
func deleteContact(c *gin.Context) {
	id := c.Param("id")

	conn, err := getConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if conn != nil {
			if err := conn.Close(context.Background()); err != nil {
				log.Println("Error closing the database connection:", err)
			}
		}
	}()

	_, err = conn.Exec(context.TODO(),
		"DELETE FROM contacts WHERE id = $1",
		id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted"})
}
func getConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.TODO(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
