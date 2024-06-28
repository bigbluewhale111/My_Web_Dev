package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type JWTClaim struct {
	Id uint32 `json:"id"`
	jwt.RegisteredClaims
}

type User struct {
	gorm.Model
	ID       uint32 `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
}

type JWTAccessToken struct {
	AccessToken string `json:"access_token"`
	jwt.RegisteredClaims
}
