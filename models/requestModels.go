package models

type RequestRegister struct {
	MbEmail           string `json:"email" form:"email" `
	Mbusername        string `json:"username" form:"username" `
	GeneratedPassword string
}

type RequestLogin struct {
	MbEmail    string `json:"email" form:"email" `
	MbPassword string `json:"password" form:"password" `
}
