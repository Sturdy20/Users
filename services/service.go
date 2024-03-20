package services

import (
	"users/models"
	"users/repositories"
	"log"
)

type IService interface {
	AddLoginService(login models.RequestLogin) error
	AddRegisterService(register models.RequestRegister) error
}

type service struct {
	r repositories.IRepositorie
}

func NewService(r repositories.IRepositorie) IService {
	return &service{r: r}
}
func (s *service) AddRegisterService(register models.RequestRegister) error {
	err := s.r.AddRegisterRepositorie(register)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (s *service) AddLoginService(login models.RequestLogin) error {
	err := s.r.AddLoginRepositorie(login)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
