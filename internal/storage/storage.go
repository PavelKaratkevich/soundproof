package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	Domain "soundproof/internal/domain"

	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	db *sqlx.DB
}

func (s *PostgreSQL) RegisterUserInDB(ctx context.Context, req Domain.UserRegistrationRequest) (int, error) {
	sqlRequest := "INSERT INTO users (Username, Password, Fullname, Email) VALUES ($1, $2, $3, $4)"

	res, err := s.db.Exec(sqlRequest, req.Username, req.Password, req.FullName, req.Email)
	if err != nil {
		return 0, err
	}

	rowsAdded, _ := res.RowsAffected()

	return int(rowsAdded), nil
}

func ConnectPostgresDB() *sqlx.DB {

	// load environment variables
	DB_DRIVER := os.Getenv("DB_DRIVER")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PORT := os.Getenv("DB_PORT")
	DB_TABLE := os.Getenv("DB_TABLE")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	// connect DB
	db, err := sqlx.Open(DB_DRIVER, fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_TABLE))
	if err != nil {
		log.Fatalf("Error while opening DB: ", err)
	}

	// ping database
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while pinging the database: ", err)
	}
	return db
}

func NewPostgreSQL(conn *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: conn,
	}
}
