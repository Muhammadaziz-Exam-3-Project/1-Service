package storage

import (
	"database/sql"
	"fmt"
	"log"

	config "Github.com/LocalEats/Order-Service/internal/config"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func ConnectDB(config config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		config.DB_USER,
		config.DB_NAME,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	log.Printf("--------------------------- Connected to the database %s --------------------------------\n", config.DB_NAME)

	return db, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
