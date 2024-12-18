package storages

import (
	"context"
	"errors"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type AuthStorage struct {
	databasePool *pgxpool.Pool
}

func NewAuthStorage(dbpool *pgxpool.Pool) *AuthStorage {
	authStorage := new(AuthStorage)
	authStorage.databasePool = dbpool
	return authStorage
}

func (storage *AuthStorage) FindOneUserByEmail(email string) models.User {
	query := "SELECT id, full_name, rank, phone, email, password_hash FROM users WHERE email = $1"
	var result models.User
	err := pgxscan.Get(context.Background(), storage.databasePool, &result, query, email)

	if err = errors.Unwrap(errors.Unwrap(err)); err != nil && err != pgx.ErrNoRows {
		logrus.Errorln(err)
	}

	return result
}

func (storage *AuthStorage) CreateOneUser(userData models.Auth, passwordHash string) (models.User, error) {
	query := "INSERT INTO users(full_name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, email, full_name"
	var user models.User

	err := pgxscan.Get(context.Background(), storage.databasePool, &user, query, userData.FullName, userData.Email, passwordHash)

	if err != nil {
		logrus.Errorln(err)
		return user, err
	}

	return user, nil
}
