package repositories

import (
	"database/sql"
	"errors"
	"log"
	"users/models"
)

type IRepositorie interface {
	AddLoginRepositorie(login models.RequestLogin) error
	AddRegisterRepositorie(register models.RequestRegister) error
}

type repositorie struct {
	db *sql.DB
}

func NewRepositorie(db *sql.DB) IRepositorie {
	return &repositorie{db: db}

}
func (r *repositorie) AddRegisterRepositorie(register models.RequestRegister) error {
	var mbID string

	err := r.db.QueryRow("SELECT mb_id FROM members WHERE mb_email = $1", register.MbEmail).Scan(&mbID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("failed to query members table. Error: %v", err)
		return errors.New("failed to query members table (API addRegister)")
	}

	if err == nil {
		// หากมีอีเมลนี้อยู่แล้วในระบบ
		return errors.New("email already exists")
	}

	// หากไม่พบสมาชิก ให้สร้างสมาชิกใหม่
	err = r.db.QueryRow("INSERT INTO members (mb_username, mb_email, mb_password) VALUES ($1, $2, $3) RETURNING mb_id", register.Mbusername, register.MbEmail, register.MbPassword).Scan(&mbID)
	if err != nil {
		log.Printf("failed to insert into members. Error: %v", err)
		return errors.New("failed to insert into members table (API addRegister)")
	}

	// ค้นหา ID ของบทบาท
	var roleID string
	roleName := "User"
	err = r.db.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", roleName).Scan(&roleID)
	if err != nil {
		log.Printf("failed to query roles table. Error: %v", err)
		return errors.New("failed to query roles table (API addRegister)")
	}

	// ให้สร้างความสัมพันธ์ระหว่างสมาชิกและบทบาท
	_, err = r.db.Exec("UPDATE members SET mb_role_id = $1 WHERE mb_id = $2", roleID, mbID)
	if err != nil {
		log.Printf("failed to update member role. Error: %v", err)
		return errors.New("failed to update member role (API addRegister)")
	}

	return nil
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

	// ตรวจสอบความถูกต้องของรหัสผ่าน
	if storedPassword != login.MbPassword {
		return errors.New("the password is incorrect")
	}

	return nil
}
