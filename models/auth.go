package models

import (
	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	UserId uint
	jwt.StandardClaims
}