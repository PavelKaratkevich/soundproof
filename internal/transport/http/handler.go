package transport

import (
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

// RegisterUser		 	godoc
// @Summary      		Register a user
// @Description  		Register a user by passing a User Registration Request via the context
// @Accept       		json
// @Produce      		json
// @Success      		200  {object}  error
// @Failure      		403  {object}  error
// @Failure      		500  {object}  error
// @Router       		/auth/register/ [post].
func (h Handler) RegisterUser(c *gin.Context) {
	var newRequest domain.UserRegistrationRequest

	// Validating if all the fields are filled in
	if err := c.ShouldBindJSON(&newRequest); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RegisterUser(newRequest)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		if err.Error() == "user with this email has already been registered" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status:": "user has been registered successfully"})
	}
}

// UpdateUser		 	godoc
// @Summary      		Update user info
/* @Description  		Update user info by passing a User Update Request via the context.
Parses signature and signed string from Metamask and stores a Metamask public address into the database */
// @Accept       		json
// @Produce      		json
// @Success      		200  {object}  error
// @Failure      		401  {object}  error
// @Failure      		404  {object}  error
// @Failure      		500  {object}  error
// @Router       		/user/profile [put].
func (h Handler) UpdateUser(c *gin.Context) {
	var newRequest domain.UpdateUserProfileRequest

	// Validating if all the fields are filled in
	if err := c.ShouldBindJSON(&newRequest); err != nil {
		log.Printf("Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ifValid, _, err := h.service.CheckCredentials(newRequest.Email, newRequest.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User with this email address not found"})
			return
		} else if err.Error() == "wrong password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	if ifValid {
		if err := jwtauth.TokenValid(c.Request); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please provide a valid authentication token"})
			return
		}

		err = h.service.UpdateUser(newRequest)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status:": "user info has been updated successfully"})
	}
}

// GetUser			 	godoc
// @Summary      		Gets a user
// @Description  		Received login request via context, checks JWT token and retrieves user info
// @Accept       		json
// @Produce      		json
// @Success      		200  {object} error
// @Failure      		403  {object}  error
// @Failure      		500  {object}  error
// @Router       		/user/profile [get].
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

	resp, err := h.service.GetUserProfile(req)
	if err != nil {
		log.Println(err)
		if err.Error() == "please provide valid credentials" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	} else {
		c.JSON(http.StatusOK, resp)
		return
	}
}

// Login			 	godoc
// @Summary      		Login form
// @Description  		Login form which received login/password, generates JWT token and returns a login response (user info)
// @Accept       		json
// @Produce      		json
// @Success      		200  {object}  error
// @Failure      		403  {object}  error
// @Failure      		500  {object}  error
// @Router       		/auth/login [post].
func (h Handler) Login(c *gin.Context) {
	var req domain.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ifValid, user, err := h.service.CheckCredentials(req.Email, req.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User with this email address is not found"})
		} else if err.Error() == "wrong password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
	}

	return nil
}
