package storages

import (
	"context"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type UsersStorage struct {
	databasePool *pgxpool.Pool
}

func NewUsersStorage(pool *pgxpool.Pool) *UsersStorage {
	usersStorage := new(UsersStorage)
	usersStorage.databasePool = pool
	return usersStorage
}

func (storage *UsersStorage) FindAllUsers() []models.User {

	query := "SELECT * FROM users"
	var result []models.User
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query)

	if err != nil {
		logrus.Errorln(err)
	}

	return result
}

func (storage *UsersStorage) FindOneUserById(id int64) models.User {

	query := "SELECT * FROM users WHERE id = $1"
	var result models.User
	err := pgxscan.Get(context.Background(), storage.databasePool, &result, query, id)

	if err != nil {
		logrus.Errorln(err)
	}

	return result
}
