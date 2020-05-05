package models

import "time"

//go:generate easyjson -all
type Post struct {
	ID       int       `json:"id,omitempty"`
	Parent   int       `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int       `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Path     []int64   `json:"-"`
}

//easyjson:json
type Posts []Post
