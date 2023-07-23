package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	Domain "soundproof/internal/domain"

	"github.com/jmoiron/sqlx"
)

// PostgreSQL implements the interface Storage.
type PostgreSQL struct {
	logger *zap.Logger
	db *sqlx.DB
}

func (s *PostgreSQL) RegisterUserInDB(ctx *gin.Context, req Domain.UserRegistrationRequest) (int, error) {
	s.logger.Debug(">>>>>> Updating the database with user registration form")

	sqlRequest := "INSERT INTO public.users (firstname, password, full_name, email) VALUES ($1, $2, $3, $4)"

	res, err := s.db.Exec(sqlRequest, req.Username, req.Password, req.FullName, req.Email)
	if err != nil {
		return 0, err
	}

	rowsAdded, _ := res.RowsAffected()

	return int(rowsAdded), nil
}

func ConnectPostgresDB(logger *zap.Logger) *sqlx.DB {

	logger.Debug(">>>>>>> Connecting PostgreSQL database")

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
		log.Fatalf("Error while opening DB: %v", err)
	}

	// ping database
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while pinging the database: %v", err)
	}
	return db
}

func NewPostgreSQL(logger  *zap.Logger, conn *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: conn,
		logger: logger,
	}
}
