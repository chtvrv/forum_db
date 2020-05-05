package models

//go:generate easyjson -all
type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"thread,omitempty"`
}
