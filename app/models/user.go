package models

//go:generate easyjson -all
type User struct {
	Nickname string `json:"nickname, omitempty"`
	Fullname string `json:"fullname, omitempty"`
	About    string `json:"about, omitempty"`
	Email    string `json:"email, omitempty"`
}

//easyjson:json
type Users []User
