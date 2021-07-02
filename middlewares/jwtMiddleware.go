package middlewares

import (
	"context"
	"cryptocurrencies-api/config"
	"cryptocurrencies-api/models"
	u "cryptocurrencies-api/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func JwtValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, config.ControllersConfig.Messages["MissingToken"])
			return
		}

		if !strings.HasPrefix(tokenHeader, "Bearer ") {
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, config.ControllersConfig.Messages["Invalid Authorization field"])
			return
		}

		tokenClaims := &config.TokenClaims{}

		token, err := jwt.ParseWithClaims(
			strings.TrimPrefix(tokenHeader, "Bearer "),
			tokenClaims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(config.JwtConfig.Secret), nil
			})

		if ve, ok := err.(*jwt.ValidationError); !token.Valid && ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				w.WriteHeader(http.StatusForbidden)
				u.Respond(w, config.ControllersConfig.Messages["MalformedToken"])
				return
			case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				w.WriteHeader(http.StatusForbidden)
				u.Respond(w, config.ControllersConfig.Messages["ExpiredOrNotActiveToken"])
				return
			case ve.Errors&(jwt.ValidationErrorClaimsInvalid) != 0:
				w.WriteHeader(http.StatusForbidden)
				u.Respond(w, config.ControllersConfig.Messages["ValidationErrorClaimsInvalid"])
				return
			case ve.Errors&(jwt.ValidationErrorSignatureInvalid) != 0:
				w.WriteHeader(http.StatusForbidden)
				u.Respond(w, config.ControllersConfig.Messages["ValidationErrorSignatureInvalid"])
				return
			default:
				fmt.Println(err)
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), "UserId", tokenClaims.UserId))
		next.ServeHTTP(w, r)
	})
}

func JwtRefreshValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, config.ControllersConfig.Messages["MissingToken"])
			return
		}

		token := strings.TrimPrefix(tokenHeader, "Bearer ")

		userToken := models.TokenSchema{UserId: r.Context().Value("UserId")}
		if values, err := models.TokenStorage.SelectValues([]interface{}{&userToken}, []string{"userId"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			u.Respond(w, config.ControllersConfig.Messages["InternalServerError"])
			return
		} else if u.GetField(values, "RefreshToken") != token {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, config.ControllersConfig.Messages["NotRelevantToken"])
			return
		}

		next.ServeHTTP(w, r)
	})
}
