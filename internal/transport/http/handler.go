package transport

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net/http"
	"reflect"
	domain "soundproof/internal/domain/model"
	"soundproof/internal/domain/service"
	jwtauth "soundproof/internal/transport/middleware/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (h Handler) GetUserByItsID(c *gin.Context) {

	// retrieve id from context
	id := c.Param("id")

	// launch validator
	validate := validator.New()

	// validate if id is an integer, required, greater than 0
	errs := validate.Var(id, "required,number")
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
		return
	}

	// if valid, we need to make sure JWT token is valid
		if err := jwtauth.TokenValid(c.Request); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please provide a valid authentication token"})
			return
		}
	
	// convert id to integer to pass it to the service level as an argument
	id_int, err := strconv.Atoi(id)

	resp, err := h.service.GetByID(c, id_int)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with ID: %v not found", id)})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ifValid, user, err := h.service.CheckCredentials(c, req)
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

func NewHandler(logger *zap.Logger, s *service.UserService) *Handler {
	return &Handler{
		service: *s,
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
