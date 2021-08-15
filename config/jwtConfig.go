package config

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
)

type TokenClaims struct {
	UserId interface{}
	jwt.StandardClaims
}

type jwtConfigStruct struct {
	Secret                     string
	AccessTokenExpiresMinutes  uint
	RefreshTokenExpiresMinutes uint
}

var JwtConfig *jwtConfigStruct

func init() {
	if SECRET, ok := os.LookupEnv("JWT_SECRET"); ok {
		JwtConfig = &jwtConfigStruct{
			Secret:                     SECRET,
			AccessTokenExpiresMinutes:  60,
			RefreshTokenExpiresMinutes: 2880,
		}
	} else {
		log.Fatal("Needed ENVIRONMENT VARIABLES are absent")
	}
}
