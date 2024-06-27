package models

import (
	"time"

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
	gorm.Model
	ID        uint32 `json:"id"`
	Token     string `json:"token"`
	UserID    uint32 `json:"user_id"`
	ExpiredAt time.Time
	CreatedAt time.Time
}

type JWTClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ValidateRequest struct {
	Token string `json:"token"`
}
