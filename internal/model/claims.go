package model

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserId int
	jwt.StandardClaims
}
