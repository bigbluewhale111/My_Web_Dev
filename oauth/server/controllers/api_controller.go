package controllers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bigbluewhale111/oauth_server/models"
	"golang.org/x/crypto/bcrypt"
)

var secret = os.Getenv("SECRET")

func (c controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error parsing form"))
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	retype_password := r.FormValue("retype_password")

	if password != retype_password {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Err: Password and retype password do not match"))
		return
	}

	if result := c.DB.First(&models.User{}, "username = ?", username); result.Error == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Err: User already exists"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error hashing password"))
		return
	}

	if result := c.DB.Create(&models.User{
		Username: username,
		Password: string(hashedPassword),
	}); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

func (c controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error parsing form"))
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username, password)
	var user models.User
	if result := c.DB.First(&user, "username = ?", username); result.Error != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Err: Invalid username or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Err: Invalid username or password"))
		return
	}

	authorize_code := fmt.Sprintf("%x", md5.Sum([]byte(username+secret+time.Now().String())))
	err = c.RDB.Set(r.Context(), "authcode:"+authorize_code, username, time.Minute*10).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error saving auth code"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(authorize_code))
}

func (c controller) ValidateToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request body"))
		return
	}
	var req models.ValidateRequest
	json.Unmarshal(body, &req)

	username, err := c.RDB.Get(r.Context(), "token:"+req.Token).Result()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(username))
}

func (c controller) RevokeToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request body"))
		return
	}
	var req models.ValidateRequest
	json.Unmarshal(body, &req)

	err = c.RDB.Del(r.Context(), "token:"+req.Token).Err()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token revoked"))
}

func (c controller) SendToken(w http.ResponseWriter, r *http.Request) {

	redirect_uri := r.FormValue("redirect_uri")
	if redirect_uri != "http://localhost:7003/callback" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Hacker ?"))
		return
	}
	code := r.FormValue("code")
	userName, err := c.RDB.Get(r.Context(), "authcode:"+code).Result()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid code"))
		return
	}

	AccessToken := sha256.Sum256([]byte(userName + secret + time.Now().String()))

	tokenString := fmt.Sprintf("%x", AccessToken)

	err = c.RDB.Set(r.Context(), "token:"+tokenString, userName, time.Hour).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error saving token"))
		return
	}
	fmt.Println(tokenString)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AccessToken{
		Token:  tokenString,
		Expire: 3600,
	})

}

func ServeFileLogin(w http.ResponseWriter, r *http.Request) {
	// Serve the login.html file
	http.ServeFile(w, r, "./dist/login.html")
}
func ServeFileRegister(w http.ResponseWriter, r *http.Request) {
	// Serve the login.html file
	http.ServeFile(w, r, "./dist/register.html")
}
