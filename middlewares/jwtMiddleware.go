package middlewares

import (
	"context"
	"cryptocurrencies-api/config"
	u "cryptocurrencies-api/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var (
	MissingToken = u.Message(false, "Missing authentication Token")
	MalformedToken = u.Message(false, "Malformed authentication Token")
	InvalidToken = u.Message(false, "Invalid authentication Token")
	ExpiredOrNotActiveToken = u.Message(false, "Authentication Token is either expired or not active yet")
)

func JwtValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, MissingToken)
			return
		}
		defer fmt.Println(tokenHeader)

		tokenPart := strings.Split(tokenHeader, "Bearer ")
		if len(tokenPart) != 2{
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, MalformedToken)
			return
		}
		defer fmt.Println(tokenPart)
		tokenClaims := &config.TokenClaims{}

		jwtToken := tokenPart[1]
		token, err := jwt.ParseWithClaims(jwtToken, tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtConfig.Secret), nil
		})

		if ve, ok := err.(*jwt.ValidationError); !token.Valid && ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				w.WriteHeader(http.StatusForbidden)
				u.Respond(w, InvalidToken)
				return
			case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				w.WriteHeader(http.StatusForbidden)
				u.Respond(w, ExpiredOrNotActiveToken)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "userId", tokenClaims.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func JwtRefreshValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
