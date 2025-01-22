package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/Swiffy12/taskify/src/internals/app/storages"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage *storages.AuthStorage
}

func NewAuthService(storage *storages.AuthStorage) *AuthService {
	authService := new(AuthService)
	authService.storage = storage
	return authService
}

func (service *AuthService) Register(authModel models.Auth) (string, error) {
	existedUser := service.storage.FindOneUserByEmail(authModel.Email)
	if existedUser.Email == authModel.Email {
		return "", errors.New("данный пользователь уже существует")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authModel.Password), 10)
	if err != nil {
		return "", errors.New("пароль не должен быть слишком большим")
	}

	var createdUser models.User
	createdUser, err = service.storage.CreateOneUser(authModel, string(passwordHash))
	if err != nil {
		return "", errors.New("не удалось создать пользователя")
	}

	payload := jwt.MapClaims{
		"sub": strconv.FormatUint(createdUser.Id, 10),
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}

	tokenMap := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := tokenMap.SignedString([]byte(os.Getenv("TASKIFY_JWT_SECRET_KEY")))
	if err != nil {
		return "", errors.New("не удалось создать токен авторизации пользователя")
	}

	return token, nil
}

func (service *AuthService) Login(email string, password string) (string, error) {
	existedUser := service.storage.FindOneUserByEmail(email)
	if existedUser.Email != email {
		return "", errors.New("пользователь не найден")
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("неверные данные авторизации")
	}

	payload := jwt.MapClaims{
		"sub": strconv.FormatUint(existedUser.Id, 10),
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}

	tokenMap := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := tokenMap.SignedString([]byte(os.Getenv("TASKIFY_JWT_SECRET_KEY")))
	if err != nil {
		return "", errors.New("не удалось создать токен авторизации пользователя")
	}

	return token, nil
}
