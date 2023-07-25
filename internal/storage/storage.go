package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"soundproof/config"
	Domain "soundproof/internal/domain/model"
	domain "soundproof/internal/domain/model"

	"github.com/jmoiron/sqlx"
)

// PostgreSQL implements the interface Storage.
type PostgreSQL struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func (s *PostgreSQL) RegisterUserInDB(ctx *gin.Context, req Domain.UserRegistrationRequest) (int, error) {
	s.logger.Debug(">>>>>> Updating the database with user registration form")

	sqlRequest := "INSERT INTO public.users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)"

	res, err := s.db.Exec(sqlRequest, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return 0, err
	}

	rowsAdded, _ := res.RowsAffected()

	return int(rowsAdded), nil
}

func (s *PostgreSQL) GetUserByID(ctx *gin.Context, id int) (*domain.ProfileResponse, error) {

	var user domain.ProfileResponse

	sqlRequest := "SELECT id, first_name, last_name, email, created_at FROM public.users WHERE id = $1"

	err := s.db.Get(&user, sqlRequest, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no results found")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (s *PostgreSQL) CheckUserCredentials(ctx *gin.Context, req domain.LoginRequest) (bool, *domain.LoginResponse, error) {
	s.logger.Debug(">>>>>> Checking if credentials are valid...................")

	var user domain.User

	sqlRequest := "SELECT * FROM public.users WHERE email = $1"

	err := s.db.Get(&user, sqlRequest, req.Email)
	if err != nil {
		return false, nil, err
	}

	if req.Password != user.Password {
		return false, nil, fmt.Errorf("wrong password")
	}

	// if creds are OK, we return user info (all but password)
	if req.Email == user.Email || req.Password == user.Password {
		return true, &domain.LoginResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Created:   user.Created,
		}, nil
	} else {
		return false, nil, fmt.Errorf("unknown server error")
	}
}

func ConnectPostgresDB(logger *zap.Logger, cfg *config.Config) *sqlx.DB {

	logger.Debug(">>>>>>> Connecting PostgreSQL database")

	// connect DB
	db, err := sqlx.Open(cfg.Connection.DB_DRIVER, fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", cfg.Connection.DB_USER, cfg.Connection.DB_PASSWORD, cfg.Connection.DB_HOST, cfg.Connection.DB_PORT, cfg.Connection.DB_TABLE))
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

func NewPostgreSQL(logger *zap.Logger, conn *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db:     conn,
		logger: logger,
	}
}
