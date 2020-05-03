package models

//go:generate easyjson -all
type Forum struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

//easyjson:json
type Forums []Forum
