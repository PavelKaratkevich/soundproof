package domain

import (
	"time"

	"github.com/gin-gonic/gin"
)

// User registration part
type UserRegistrationRequest struct {
	FirstName string `json:"first_name" binding:"required,alphanum"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
}

// User update part
type UpdateUserProfileRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	SignedMessage string `json:"signed_message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// Domain user
type User struct {
	ID        int       `json:"id" binding:"required,alphanum"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Password  string    `json:"password" db:"password"`
	Email     string    `json:"email" db:"email"`
	Created   time.Time `json:"created_at" db:"created_at"`
}

// Login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login response
type LoginResponse struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email" db:"email"`
	Created      time.Time `json:"created_at" db:"created_at"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
}

// Profile request/response part
type GetProfileRequest struct {
	ID int `json:"id" db:"id" binding:"required,alphanum"`
}

type ProfileResponse struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Created   time.Time `json:"created_at" db:"created_at"`
}

// Port for storage implementation.
type Storage interface {
	RegisterUserInDB(ctx *gin.Context, req UserRegistrationRequest) error
	CheckUserCredentials(ctx *gin.Context, email, password string) (bool, *LoginResponse, error)
	GetUserProfile(ctx *gin.Context, req LoginRequest) (*ProfileResponse, error)
	UpdateUserProfile(ctx *gin.Context, address, email string) error
}

// Port for service implementation.
type Service interface {
	RegisterUser(c *gin.Context, req UserRegistrationRequest) error
	GetUserProfile(c *gin.Context, req LoginRequest) (*ProfileResponse, error)
	CheckCredentials(c *gin.Context, email, password string) (bool, *LoginResponse, error)
	UpdateUser(c *gin.Context, req UpdateUserProfileRequest) error
}
