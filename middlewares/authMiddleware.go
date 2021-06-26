package middlewares

import (
	"context"
	"cryptocurrencies-api/models"
	u "cryptocurrencies-api/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := make(map[string] interface{})
		tokenHeader := r.Header.Get("Authorization") //Получение токена

		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, u.Message(false, "Missing auth token"))
			return
		}

		matched, err := regexp.Match(`Bearer [\d\w\-\.]+`, []byte(tokenHeader))
		if !matched {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, u.Message(false, "Missing auth token"))
			return
		}

		tokenPart := strings.Split(tokenHeader, " ")[1]
		tk := &models.TokenClaims{}
		defer fmt.Println(tokenHeader)
		defer fmt.Println(tokenPart)


		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		defer fmt.Printf("User %v", tk.UserId)


		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
