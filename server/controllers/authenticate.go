package controllers

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"net/http"

	"github.com/bigbluewhale111/rest_api/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserKey string

func (c controller) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get Cookie token
		tokenString, err := r.Cookie("token")
		if err != nil {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// jwt
		token, err := jwt.ParseWithClaims(tokenString.Value, &models.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claim := token.Claims.(*models.JWTClaim)
		if !token.Valid {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		Md5TokenString := fmt.Sprintf("%x", md5.Sum([]byte(tokenString.Value)))
		fmt.Println(Md5TokenString)
		login, err := c.RDB.Get(context.Background(), "session:"+Md5TokenString).Result()
		if err != nil {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if login != "true" {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println(login)

		// Add user to context
		ctx := context.WithValue(r.Context(), UserKey("user"), *claim)
		ctx = context.WithValue(ctx, UserKey("token"), Md5TokenString)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
