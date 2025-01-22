package services

import (
	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/Swiffy12/taskify/src/internals/app/storages"
)

type UsersService struct {
	storage *storages.UsersStorage
}

func NewUsersService(storage *storages.UsersStorage) *UsersService {
	usersService := new(UsersService)
	usersService.storage = storage
	return usersService
}

func (service *UsersService) GetUsersWithFilter(queryParams models.GetUsersRequestDTO) ([]models.GetUserResponseDTO, error) {
	return service.storage.GetUsersWithFilter(queryParams)
}

func (service *UsersService) GetOneUser(id int64) (*models.GetUserResponseDTO, error) {
	return service.storage.FindOneUserById(id)
}
