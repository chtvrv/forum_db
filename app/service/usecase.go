package service

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Usecase interface {
	GetStatus() (*models.Status, error)
	ClearDB() error
}
