package models

type RequestRegister struct {
	MbEmail    string `json:"mb_email" form:"mb_email" `
	Mbusername string `json:"mb_username" form:"mb_username" `
	MbPassword string `json:"mb_password" form:"mb_password" `
	MbroleID   string `json:"mb_role_id" form:"mb_role_id" `
}

type RequestLogin struct {
	MbEmail    string `json:"email" form:"email" `
	MbPassword string `json:"password" form:"password" `
}
