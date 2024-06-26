package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type JWTClaim struct {
	Username string `json:"username"`
	Id       uint32 `json:"id"`
	jwt.RegisteredClaims
}

type OauthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUser struct {
	Login string `json:"login"`
	Id    uint32 `json:"id"`
}

type TokenResponse struct {
	Id        uint32     `json:"id"`
	Auth_User GithubUser `json:"user"`
}

type User struct {
	gorm.Model
	ID          uint32 `json:"id" gorm:"primaryKey"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
	GithubId    uint32 `json:"github_id"`
}
