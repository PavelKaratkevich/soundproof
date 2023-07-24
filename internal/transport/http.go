package transport

import (
	"log"
	"net/http"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/service"
	jwtauth "soundproof/internal/transport/middleware/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service service.UserService
}

func (h Handler) RegisterUser(c *gin.Context) {
	var newRequest domain.UserRegistrationRequest

	// Validating if all the fields are filled in
	if err := c.ShouldBindJSON(&newRequest); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.RegisterUser(c, newRequest)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (h Handler) Login(c *gin.Context) {
	var req domain.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	} 

	ifValid, user, err := h.service.CheckCredentials(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} 

	if ifValid {
		ts, err := jwtauth.CreateToken()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		user.AccessToken = ts.AccessToken
		user.RefreshToken = ts.RefreshToken

		c.JSON(http.StatusOK, user)
	}	
	
}

func NewHandler(logger *zap.Logger, s *service.UserService) *Handler {
	return &Handler{
		service: *s,
		logger:  logger,
	}
}
