package domain

import (
	"time"

	"github.com/gin-gonic/gin"
)

type UserRegistrationRequest struct {
	FirstName string `json:"first_name" binding:"required,alphanum"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Password  string    `json:"password" db:"password"`
	Email     string    `json:"email" db:"email"`
	Created   time.Time `json:"created_at" db:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserRegistrationResponse struct {
	userID int `json:"user_id" binding:"required"`
}

// Port for database implementation.
type Storage interface {
	RegisterUserInDB(ctx *gin.Context, req UserRegistrationRequest) (int, error)
	CheckUserCredentials(ctx *gin.Context, req LoginRequest) (bool, error)
}
