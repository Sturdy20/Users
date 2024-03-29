package models

type RequestRegister struct {
	Mbusername        string `json:"mb_username" form:"mb_username" `
	MbEmail           string `json:"mb_email" form:"mb_email" `
	GeneratedPassword string
}

type RequestLogin struct {
	MbEmail    string `json:"email" form:"email" `
	MbPassword string `json:"password" form:"password" `
}
