package controllers

import (
	"context"
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
		claims, ok := token.Claims.(*models.JWTClaim)
		if !ok || !token.Valid {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// temporary get to the database, we will setup redis and read cache which is much faster
		var user models.User
		if result := c.DB.First(&user, "id = ? AND username = ?", claims.Id, claims.Username); result.Error != nil {
			log.Println("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("Authorized")
		fmt.Println(user)

		// Add user to context
		ctx := context.WithValue(r.Context(), UserKey("user"), user)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
