package usecase

import (
	"github.com/chtvrv/forum_db/app/service"
	"github.com/labstack/gommon/log"

	"github.com/chtvrv/forum_db/app/models"
)

type ServiceUsecase struct {
	serviceRepo service.Repository
}

func CreateUsecase(serviceRepo_ service.Repository) service.Usecase {
	return &ServiceUsecase{serviceRepo: serviceRepo_}
}

func (usecase *ServiceUsecase) GetStatus() (*models.Status, error) {
	status, err := usecase.serviceRepo.GetStatus()
	if err != nil {
		log.Error(nil)
		return nil, err
	}
	return status, nil
}

func (usecase *ServiceUsecase) ClearDB() error {
	err := usecase.serviceRepo.ClearDB()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
