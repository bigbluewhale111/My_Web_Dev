package controllers

import (
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

	// I want to generate a random hash length 64 bytes

	AccessToken := sha256.Sum256([]byte(password + secret + time.Now().String()))

	tokenString := username + "_" + fmt.Sprintf("%x", AccessToken)

	var TokenNumbers int64
	c.DB.Model(&models.AccessToken{}).Where("user_id = ?", user.ID).Count(&TokenNumbers)

	if TokenNumbers > 10 {
		if result := c.DB.Model(&models.AccessToken{}).Where("user_id = ? ", user.ID).Order("created_at asc").Limit(1).Updates(models.AccessToken{
			Token:     tokenString,
			ExpiredAt: time.Now().Add(time.Hour * 6),
			CreatedAt: time.Now()}); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error saving token"))
			return
		}
	} else {
		if result := c.DB.Create(&models.AccessToken{
			Token:     tokenString,
			UserID:    user.ID,
			ExpiredAt: time.Now().Add(time.Hour * 6),
		}); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error saving token"))
			return
		}
	}
	fmt.Println(tokenString)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
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

	var AccessToken models.AccessToken
	if result := c.DB.First(&AccessToken, "token = ?", req.Token); result.Error != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid token"))
		return
	}

	if time.Now().After(AccessToken.ExpiredAt) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Token expired"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Valid token"))
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

	var AccessToken models.AccessToken
	if result := c.DB.First(&AccessToken, "token = ?", req.Token); result.Error != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid token"))
		return
	}

	if result := c.DB.Delete(&AccessToken); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error deleting token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token revoked"))
}

func ServeFileLogin(w http.ResponseWriter, r *http.Request) {
	// Serve the login.html file
	http.ServeFile(w, r, "./dist/login.html")
}
func ServeFileRegister(w http.ResponseWriter, r *http.Request) {
	// Serve the login.html file
	http.ServeFile(w, r, "./dist/register.html")
}
