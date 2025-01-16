package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Swiffy12/taskify/src/internals/app/handlers"
	"github.com/Swiffy12/taskify/src/internals/constants"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func CheckResolution(whitelist []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if isWhitelisted(r.URL.Path, whitelist) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.WrapErrorUnauthorized(w, errors.New("пропущен заголовок авторизации"))
				return
			}

			bearer := strings.Split(authHeader, " ")
			if len(bearer) != 2 || bearer[0] != "Bearer" {
				handlers.WrapErrorUnauthorized(w, errors.New("неверный формат заголовка авторизации"))
				return
			}

			tokenString := bearer[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}

				return []byte(os.Getenv("TASKIFY_JWT_SECRET_KEY")), nil
			})

			if err != nil || !token.Valid {
				if err != nil {
					logrus.Errorln(err)
				}
				handlers.WrapErrorUnauthorized(w, errors.New("токен не действителен"))
				return
			}

			userId, err := token.Claims.GetSubject()
			if err != nil {
				handlers.WrapErrorUnauthorized(w, errors.New("недопустимые данные токена"))
			}

			ctx := context.WithValue(r.Context(), constants.UserIdKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isWhitelisted(path string, whitelist []string) bool {
	for _, p := range whitelist {
		if p == path {
			return true
		}
	}
	return false
}
