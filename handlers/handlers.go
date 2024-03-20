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
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "OK", "message": "Register Successfully"})
}



func (h *handler) AddLoginHandler(c *gin.Context) {
	var login models.RequestLogin

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	err := h.s.AddLoginService(login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "login Successfully"})
}

