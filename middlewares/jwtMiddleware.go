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
		fmt.Println("JwtValidation")
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			config.ControllersConfig.Responses["MissingToken"](w)
			return
		}

		if !strings.HasPrefix(tokenHeader, "Bearer ") {
			config.ControllersConfig.Responses["InvalidToken"](w)
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
				config.ControllersConfig.Responses["InvalidToken"](w)
				return
			case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				config.ControllersConfig.Responses["ExpiredOrNotActiveToken"](w)
				return
			case ve.Errors&(jwt.ValidationErrorClaimsInvalid) != 0:
				config.ControllersConfig.Responses["InvalidToken"](w)
				return
			case ve.Errors&(jwt.ValidationErrorSignatureInvalid) != 0:
				config.ControllersConfig.Responses["InvalidToken"](w)
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
		fmt.Println("JwtRefreshValidation")
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			config.ControllersConfig.Responses["MissingToken"](w)
			return
		}

		token := strings.TrimPrefix(tokenHeader, "Bearer ")

		userToken := models.TokenSchema{UserId: r.Context().Value("UserId")}
		if values, err := models.TokenStorage.SelectValues([]interface{}{&userToken}, []string{"userId"}); err != nil {
			config.ControllersConfig.Responses["InternalServerError"](w)
			return
		} else if u.GetField(values, "RefreshToken") != token {
			config.ControllersConfig.Responses["NotRelevantToken"](w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
