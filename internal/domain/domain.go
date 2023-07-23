package domain

import "github.com/gin-gonic/gin"

type UserRegistrationRequest struct {
	FirstName string `json:"first_name" binding:"required,alphanum"`
	LastName string `json:"last_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type UserRegistrationResponse struct {
	userID int `json:"user_id" binding:"required"`
}

// Port for database implementation.
type Storage interface {
	RegisterUserInDB(ctx *gin.Context, req UserRegistrationRequest) (int, error)
}
