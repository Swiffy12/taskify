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

func (service *UsersService) GetAllUsers() []models.User {

	return service.storage.FindAll()
}
