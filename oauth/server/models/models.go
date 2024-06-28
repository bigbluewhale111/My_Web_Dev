package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint32 `json:"id" gorm:"primary_key;auto_increment"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password"`
}

type AccessToken struct {
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
}

type JWTClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ValidateRequest struct {
	Token string `json:"token"`
}
