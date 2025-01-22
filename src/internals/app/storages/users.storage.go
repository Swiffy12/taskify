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

func (storage *UsersStorage) GetUsersWithFilter(queryParams models.GetUsersRequestDTO) ([]models.GetUserResponseDTO, error) {

	query := "SELECT id, full_name, rank, phone, email FROM users WHERE 1=1"
	var result []models.GetUserResponseDTO
	args := make([]any, 0)
	placeholderNumber := 1

	if queryParams.FullName != "" {
		query += fmt.Sprintf(" AND full_name ILIKE $%d", placeholderNumber)
		args = append(args, fmt.Sprintf("%%%s%%", queryParams.FullName))
		placeholderNumber++
	}

	if queryParams.Rank != "" {
		query += fmt.Sprintf(" AND rank ILIKE $%d", placeholderNumber)
		args = append(args, fmt.Sprintf("%%%s%%", queryParams.Rank))
		placeholderNumber++
	}

	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, args...)

	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return result, nil
}

func (storage *UsersStorage) FindOneUserById(id int64) (*models.GetUserResponseDTO, error) {

	query := "SELECT id, full_name, rank, phone, email FROM users WHERE id = $1"
	var result models.GetUserResponseDTO
	err := pgxscan.Get(context.Background(), storage.databasePool, &result, query, id)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logrus.Errorln(err)
		}
		return nil, err
	}

	return &result, nil
}
