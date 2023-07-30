package transport

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"net/http"
	"reflect"
	domain "soundproof/internal/domain/model"
	jwtauth "soundproof/internal/transport/middleware/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service domain.Service
}

func (h Handler) RegisterUser(c *gin.Context) {
	var newRequest domain.UserRegistrationRequest

	// Validating if all the fields are filled in
	if err := c.ShouldBindJSON(&newRequest); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RegisterUser(c, newRequest)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status:": "user has been registered successfully"})
		return
	}
}

func (h Handler) UpdateUser(c *gin.Context) {
	var newRequest domain.UpdateUserProfileRequest

	// Validating if all the fields are filled in
	if err := c.ShouldBindJSON(&newRequest); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ifValid, _, err := h.service.CheckCredentials(c, newRequest.Email, newRequest.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User with this email address not found"})
			return
		} else if err.Error() == "wrong password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if ifValid {
		if err := jwtauth.TokenValid(c.Request); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please provide a valid authentication token"})
			return
		}

		err = h.service.UpdateUser(c, newRequest)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"status:": "user info has been updated successfully"})
			return
		}
	}

}

func (h Handler) GetUser(c *gin.Context) {

	var req domain.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if valid, we need to make sure JWT token is valid
	if err := jwtauth.TokenValid(c.Request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Please provide a valid authentication token"})
		return
	}

	//// checkCredentials /////


	resp, err := h.service.GetUserProfile(c, req)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (h Handler) Login(c *gin.Context) {
	var req domain.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ifValid, user, err := h.service.CheckCredentials(c, req.Email, req.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User with this email address not found"})
		} else if err.Error() == "wrong password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

func NewHandler(logger *zap.Logger, s domain.Service) *Handler {
	return &Handler{
		service: s,
		logger:  logger,
	}
}

func ValidateValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}
