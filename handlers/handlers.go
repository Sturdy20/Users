package handlers

import (
	"net/http"
	"users/models"
	"users/services"

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
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Please provide the necessary information", "code": http.StatusBadRequest, "details": err.Error()})
		return
	}
	registerResp, err := h.s.AddRegisterService(register)
	if err != nil {
		switch err.Error() {
		case "failed to insert into members table (API Register)":
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Registration failed. Please try again later.", "code": http.StatusInternalServerError, "details": err.Error()})
		case "email already exists in the system":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Sorry, this email address is already in use for registration.", "code": http.StatusBadRequest, "details": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error(), "code": http.StatusInternalServerError, "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "OK", "message": "Registration completed successfully", "code": http.StatusCreated, "data": registerResp})
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
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "User not found. Please check your email and try again.", "code": http.StatusNotFound, "details": err.Error()})
		case "incorrect password":
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Incorrect password. Please try again.", "code": http.StatusUnauthorized, "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "An error occurred. Please try again later.", "code": http.StatusInternalServerError, "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login successful!"})
}
