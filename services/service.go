package services

import (
	"fmt"
	"log"
	"users/models"
	"users/pkg/utility/generate"
	"users/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IService interface {
	AddLoginService(login models.RequestLogin) error
	AddRegisterService(register models.RequestRegister) (models.RegisterResponses, error)
}

type service struct {
	r repositories.IRepositorie
}

func NewService(r repositories.IRepositorie) IService {
	return &service{r: r}
}
func (s *service) AddRegisterService(register models.RequestRegister) (models.RegisterResponses, error) {
	generatedPassword, err := generate.GenerateRandomPassword(12)
	if err != nil {
		log.Println("failed to generate random password:", err)
		return models.RegisterResponses{}, err
	}
	fmt.Println("Generated Password:", generatedPassword)

	// เข้ารหัสรหัสผ่านที่สร้างขึ้น
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password. Error: %v\n", err)
		return models.RegisterResponses{}, err
	}

	// แปลงรหัสผ่านที่เข้ารหัสแล้วเป็น string
	hashedPasswordString := string(hashedPassword)

	// กำหนดรหัสผ่านที่เข้ารหัสแล้วให้กับ request register
	register.GeneratedPassword = hashedPasswordString

	// เรียกใช้งาน AddRegisterRepositorie ใน Repositorie เพื่อเพิ่มข้อมูลลงในฐานข้อมูล
	registerResp, err := s.r.AddRegisterRepositorie(register)
	if err != nil {
		log.Println(err.Error())
		return models.RegisterResponses{}, err
	}

	return registerResp, nil
}


func (s *service) AddLoginService(login models.RequestLogin) error {
	err := s.r.AddLoginRepositorie(login)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
