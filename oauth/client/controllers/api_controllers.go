package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bigbluewhale111/oauth_client/models"
	"github.com/golang-jwt/jwt/v5"
)

var ClientID = os.Getenv("CLIENT_ID")
var ClientSecret = os.Getenv("CLIENT_SECRET")
var OAUTH_SERVER_URL = os.Getenv("OAUTH_SERVER_URL")
var JWTSecret = os.Getenv("JWTSECRET")
var rest_api = os.Getenv("REST_API")

func HandlingCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the callback from the server
	code := r.URL.Query().Get("code")
	fmt.Println("Code: ", code)
	data := url.Values{}
	data.Add("code", code)
	data.Add("redirect_uri", "http://localhost:7003/callback")
	// make request to oauth server to get access token
	req, err := http.NewRequest("POST", OAUTH_SERVER_URL+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(ClientID, ClientSecret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Oauth client configured incorrectly"))
		return
	}

	defer resp.Body.Close()
	var access_token models.AccessToken
	err = json.NewDecoder(resp.Body).Decode(&access_token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// JWT the access token
	expiredTime := time.Now().Add(time.Duration(access_token.Expire) * time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTAccessToken{
		AccessToken: access_token.Token,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	})
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	req, err = http.NewRequest("GET", rest_api+"/auth/authorize", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	resp, err = client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to authorize"))
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(body)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	sessionToken := r.URL.Query().Get("session")
	data := url.Values{}
	data.Add("session", sessionToken)

	req, err := http.NewRequest("POST", rest_api+"/auth/logout", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error in creating request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in sending request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to logout"))
		return
	}

	token, err := jwt.ParseWithClaims(sessionToken, &models.JWTAccessToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		fmt.Println("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claim := token.Claims.(*models.JWTAccessToken)
	if !token.Valid {
		fmt.Println("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Revoke the token
	payload, err := json.Marshal(models.ValidateRequest{
		Token: claim.AccessToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	req, err = http.NewRequest("POST", OAUTH_SERVER_URL+"/oauth/revoke", strings.NewReader(string(payload)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(ClientID, ClientSecret)
	resp, err = client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to revoke token"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successfully"))
}