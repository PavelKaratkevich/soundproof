package domain

import "github.com/gin-gonic/gin"

type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserRegistrationResponse struct {
	UserID int `json:"user_id" binding:"required"`
}

// Port for database implementation
type Storage interface {
	RegisterUserInDB(ctx *gin.Context, req UserRegistrationRequest) (int, error)
}
