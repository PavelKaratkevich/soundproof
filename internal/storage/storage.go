package storage

import (
	"database/sql"
	"fmt"

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

func (s *PostgreSQL) RegisterUserInDB(ctx *gin.Context, req Domain.UserRegistrationRequest) error {

	// check if no user with the same email is found
	err := s.checkForExisingUsers(req)
	if err != nil {
		s.logger.Debug(fmt.Sprintf("Error while checking duplicated users in the database: %v", err.Error()))
		return err
	}

	// insert data into the database
	sqlRequest := "INSERT INTO public.users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)"
	_, err = s.db.Exec(sqlRequest, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		s.logger.Debug(fmt.Sprintf("Error while inserting user data into the database: %v", err.Error()))
		return err
	}
	return nil
}

func (s *PostgreSQL) GetUserProfile(ctx *gin.Context, req domain.LoginRequest) (*domain.ProfileResponse, error) {

	var user domain.ProfileResponse

	sqlRequest := "SELECT id, first_name, last_name, email, created_at FROM public.users WHERE email = $1 AND password = $2"

	err := s.db.Get(&user, sqlRequest, req.Email, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("please provide valid credentials")
		} else {
			s.logger.Debug(fmt.Sprintf("Error while retrieving user data from the database: %v", err.Error()))
			return nil, err
		}
	}

	return &user, nil
}

func (s *PostgreSQL) CheckUserCredentials(ctx *gin.Context, req domain.LoginRequest) (bool, *domain.LoginResponse, error) {

	var user domain.User

	sqlRequest := "SELECT id, first_name, last_name, password, email, created_at FROM public.users WHERE email = $1"

	err := s.db.Get(&user, sqlRequest, req.Email)
	if err != nil {
		s.logger.Debug(fmt.Sprintf("Error while retrieving user data from the database: %v", err.Error()))
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

	// connect DB
	db, err := sqlx.Open(cfg.Connection.DB_DRIVER, fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", cfg.Connection.DB_USER, cfg.Connection.DB_PASSWORD, cfg.Connection.DB_HOST, cfg.Connection.DB_PORT, cfg.Connection.DB_TABLE))
	if err != nil {
		logger.Debug(fmt.Sprintf("Error while opening DB: %v", err))
	}

	// ping database
	err = db.Ping()
	if err != nil {
		logger.Debug(fmt.Sprintf("Error while pinging the database: %v", err))
	}
	return db
}

func NewPostgreSQL(logger *zap.Logger, conn *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db:     conn,
		logger: logger,
	}
}

// checkForExisingUsers checks if no users with the same email address already exists
func (s *PostgreSQL) checkForExisingUsers(req Domain.UserRegistrationRequest) error {
	sqlRequestForDuplicates := "SELECT id FROM public.users WHERE email = $1"

	var id string

	err := s.db.Get(&id, sqlRequestForDuplicates, req.Email)

	if err != sql.ErrNoRows {
		return fmt.Errorf("user with this email has already been registered")
	}
	return nil
}

func (s *PostgreSQL) UpdateUserProfile(ctx *gin.Context, address, email string) error {
	sqlRequest := "UPDATE public.users SET metamask_address=$1 WHERE email = $2"

	// check for credentials so that each user can update only his/her own records

	_, err := s.db.Exec(sqlRequest, address, email)
	if err != nil {
		s.logger.Debug(fmt.Sprintf("Error while updating user profile: %v", err))
		return err
	}

	return nil
}
