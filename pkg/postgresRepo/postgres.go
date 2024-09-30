package postgresrepo

import (
	"database/sql"
	"fmt"
	"log"
	"musPlayer/models"

	_ "github.com/lib/pq"
)

func NewPostgresDb(dbConfig models.DatabaseConfig) (*sql.DB, error) {
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	log.Println("Opening database connection...")
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Println("Pinging database...")
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
