package storages

import (
	"context"
	"errors"
	"fmt"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
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

func (storage *UsersStorage) FindUsersWithFilter(fullname string, rank string) []models.User {

	query := "SELECT * FROM users WHERE 1=1"
	var result []models.User
	args := make([]interface{}, 0)
	placeholderNumber := 1

	if fullname != "" {
		query += fmt.Sprintf(" AND full_name ILIKE $%d", placeholderNumber)
		args = append(args, fmt.Sprintf("%%%s%%", fullname))
		placeholderNumber++
	}

	if rank != "" {
		query += fmt.Sprintf(" AND rank ILIKE $%d", placeholderNumber)
		args = append(args, fmt.Sprintf("%%%s%%", rank))
		placeholderNumber++
	}

	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, args...)

	if err != nil {
		logrus.Errorln(err)
	}

	return result
}

func (storage *UsersStorage) FindOneUserById(id int64) (models.User, error) {

	query := "SELECT * FROM users WHERE id = $1"
	var result models.User
	err := pgxscan.Get(context.Background(), storage.databasePool, &result, query, id)

	if err = errors.Unwrap(errors.Unwrap(err)); err != nil { // Дублирование кода
		if err == pgx.ErrNoRows {
			return result, errors.New("не удалось найти данного пользователя")
		}
		logrus.Errorln(err)
	}

	return result, nil
}
