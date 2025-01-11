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

func (service *UsersService) GetAllUsers(fullname string, rank string) []models.User {
	return service.storage.FindUsersWithFilter(fullname, rank)
}

func (service *UsersService) GetOneUser(id int64) models.User {
	return service.storage.FindOneUserById(id)
}
