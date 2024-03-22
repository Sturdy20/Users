package handlers

import (
	"users/models"
	"users/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	AddLoginHandler(c *gin.Context)
	AddRegisteHandler(c *gin.Context)


}

type handler struct {
	s services.IService
}

func NewHandler(s services.IService) IHandler {
	return &handler{s: s}
}

func (h *handler) AddRegisteHandler(c *gin.Context) {
	var register models.RequestRegister

	if err := c.ShouldBind(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	err := h.s.AddRegisterService(register)
	if err != nil {
		switch err.Error() {
		case "Failed to query members table":
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to register. Please try again later."})
		case "Email already exists":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Email already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "An error occurred. Please try again later."})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "OK", "message": "Register successfully"})
}




func (h *handler) AddLoginHandler(c *gin.Context) {
	var login models.RequestLogin

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	err := h.s.AddLoginService(login)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "User not found. Please check your email and try again."})
		case "incorrect password":
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Incorrect password. Please try again."})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "An error occurred. Please try again later."})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login successful!"})
}

