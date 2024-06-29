package controllers

import (
	"crypto/subtle"
	"net/http"
	"os"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ClientID = os.Getenv("CLIENT_ID")
		var ClientSecret = os.Getenv("CLIENT_SECRET")
		username, password, ok := r.BasicAuth()
		notauthUsername := subtle.ConstantTimeCompare([]byte(username), []byte(ClientID)) != 1
		notauthPassword := subtle.ConstantTimeCompare([]byte(password), []byte(ClientSecret)) != 1
		if !ok || notauthUsername || notauthPassword {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
