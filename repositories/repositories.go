package repositories

import (
	"database/sql"
	"errors"
	"log"
	"users/models"

	"golang.org/x/crypto/bcrypt"
)

type IRepositorie interface {
	AddLoginRepositorie(login models.RequestLogin) error
	AddRegisterRepositorie(register models.RequestRegister) (models.RegisterResponses, error)
}

type repositorie struct {
	db *sql.DB
}

func NewRepositorie(db *sql.DB) IRepositorie {
	return &repositorie{db: db}

}
func (r *repositorie) AddRegisterRepositorie(register models.RequestRegister) (models.RegisterResponses, error) {
	var mbID string

	err := r.db.QueryRow("SELECT mb_id FROM members WHERE mb_email = $1", register.MbEmail).Scan(&mbID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("failed to query members table. Error: %v", err)
		return models.RegisterResponses{}, errors.New("failed to query members table (API Register)")
	}

	if err == nil {
		return models.RegisterResponses{}, errors.New("email already exists in the system")
	}
	// หากไม่พบสมาชิก ให้สร้างสมาชิกใหม่
	err = r.db.QueryRow("INSERT INTO members (mb_username, mb_email, mb_password) VALUES ($1, $2, $3) RETURNING mb_id", register.Mbusername, register.MbEmail, register.GeneratedPassword).Scan(&mbID)
	if err != nil {
		log.Printf("failed to insert into members. Error: %v", err)
		return models.RegisterResponses{}, errors.New("failed to insert into members table (API Register)")
	}

	// ค้นหา ID ของบทบาท
	var roleID string
	roleName := "User"
	err = r.db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", roleName).Scan(&roleID)
	if err != nil {
		log.Printf("failed to query roles table. Error: %v", err)
		return models.RegisterResponses{}, errors.New("failed to query roles table (API Register)")
	}

	// ให้สร้างความสัมพันธ์ระหว่างสมาชิกและบทบาท
	_, err = r.db.Exec("UPDATE members SET mb_role_id = $1 WHERE mb_id = $2", roleID, mbID)
	if err != nil {
		log.Printf("failed to update member role. Error: %v", err)
		return models.RegisterResponses{}, errors.New("failed to update member role (API Register)")
	}
	
	// Get register response
	var registerResp models.RegisterResponses
	err = r.db.QueryRow("SELECT mb_id, mb_username, mb_email, mb_role_id FROM members WHERE mb_id = $1", mbID).Scan(&registerResp.MbID, &registerResp.MbUsername, &registerResp.MbEmail, &registerResp.RoleID)
	if err != nil {
		log.Printf("failed to get register response from members table. Error: %v", err)
		return models.RegisterResponses{}, errors.New("failed to get register response")
	}

	return registerResp, nil
}


func (r *repositorie) AddLoginRepositorie(login models.RequestLogin) error {
	var mbID string
	var storedPassword string

	err := r.db.QueryRow("SELECT mb_id, mb_password FROM members WHERE mb_email = $1", login.MbEmail).Scan(&mbID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found: %v", err)
			return errors.New("user not found")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(login.MbPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("the password is incorrect")
		}
		log.Printf("failed to compare hashed passwords. Error: %v\n", err)
		return err
	}
	

	return nil
}
