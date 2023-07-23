package transport

import (
	"log"
	"net/http"
	"soundproof/internal/domain/model"
	"soundproof/internal/domain/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handler struct {
	logger  *zap.Logger
	service service.UserService
}

func (h handler) RegisterUser(c *gin.Context) {
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

func NewHandler(logger  *zap.Logger, s *service.UserService) *handler {
	return &handler{
		service: *s,
		logger: logger,
	}
}
