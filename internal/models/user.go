package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Password string `json:",omitempty"`
	Username string
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}
