package models

import "github.com/golang-jwt/jwt/v5"

type AccessToken struct {
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
}

type JWTAccessToken struct {
	AccessToken string `json:"access_token"`
	jwt.RegisteredClaims
}

type ValidateRequest struct {
	Token string `json:"token"`
}
