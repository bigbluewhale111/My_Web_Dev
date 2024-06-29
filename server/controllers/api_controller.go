package controllers

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bigbluewhale111/rest_api/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)


func handleJWT(id uint32) (string, time.Time, error) {
	expiredTime := time.Now().Add(time.Hour * 6)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTClaim{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))
	return tokenString, expiredTime, err
}

func CachingUser(userId uint32, tokenString string, RDB *redis.Client, expiredTime time.Time) error {
	Md5TokenString := fmt.Sprintf("%x", md5.Sum([]byte(tokenString)))
	err := RDB.Set(context.Background(), "session:"+Md5TokenString, strconv.Itoa(int(userId)), time.Until(expiredTime)).Err()
	if err != nil {
		fmt.Println("Error storing token to cache:" + err.Error())
		return err
	}
	return nil
}

func (c controller) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authorId := r.Context().Value(UserKey("user")).(models.JWTClaim).Id

	var tasks []models.Task
	if result := c.DB.Where("author_id = ?", authorId).Order("id ASC").Find(&tasks); result.Error != nil {
		fmt.Println("Error fetching tasks:" + result.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (c controller) AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error reading body:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	var newtask models.NewTask
	json.Unmarshal(body, &newtask)
	if newtask.Name == "" || newtask.Description == "" || newtask.Status > 4 || newtask.Status < 1 || newtask.DueDate == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Field"))
		return
	}

	authorId := r.Context().Value(UserKey("user")).(models.JWTClaim).Id

	if result := c.DB.Create(&models.Task{
		Name:        newtask.Name,
		Description: newtask.Description,
		Status:      newtask.Status,
		DueDate:     newtask.DueDate,
		AuthorID:    authorId,
	}); result.Error != nil {
		fmt.Println("Error creating task:" + result.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Created task successfully")
}

func (c controller) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	var id, _ = strconv.Atoi(vars["id"])

	authorId := r.Context().Value(UserKey("user")).(models.JWTClaim).Id

	var task models.Task
	if result := c.DB.First(&task, "id = ? AND author_id = ?", id, authorId); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (c controller) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	var id, _ = strconv.Atoi(vars["id"])

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error reading body:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	var UpdateTask models.NewTask
	json.Unmarshal(body, &UpdateTask)
	if UpdateTask.Status > 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Field"))
		return
	}

	authorId := r.Context().Value(UserKey("user")).(models.JWTClaim).Id

	var task models.Task
	if result := c.DB.First(&task, "id = ? AND author_id = ?", id, authorId); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}
	if UpdateTask.Name != "" {
		task.Name = UpdateTask.Name
	}
	if UpdateTask.Description != "" {
		task.Description = UpdateTask.Description
	}
	if UpdateTask.Status != 0 {
		task.Status = UpdateTask.Status
	}
	if UpdateTask.DueDate != 0 {
		task.DueDate = UpdateTask.DueDate
	}

	if result := c.DB.Save(&task); result.Error != nil {
		fmt.Println("Error updating task:" + result.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Task updated successfully")
}

func (c controller) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	var id, _ = strconv.Atoi(vars["id"])

	authorId := r.Context().Value(UserKey("user")).(models.JWTClaim).Id

	var task models.Task
	if result := c.DB.First(&task, "id = ? AND author_id = ?", id, authorId); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	if result := c.DB.Delete(&task, "id = ? AND author_id = ?", id, authorId); result.Error != nil {
		fmt.Println("Error deleting task:" + result.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Task deleted successfully")
}

func (c controller) Authorize(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	AuthorizationHeader := r.Header.Get("Authorization")
	if AuthorizationHeader == "" {
		fmt.Println("No Authorization Header")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	tokenString := strings.Split(AuthorizationHeader, " ")[1]
	var AccessTokenClaim models.JWTAccessToken
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTSECRET")), nil
	})
	if err != nil || token == nil || !token.Valid {
		fmt.Println("Error parsing token:" + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	claimToken := token.Claims.(*models.JWTAccessToken)
	url := os.Getenv("OAUTH_URL") + "/validate"
	payload, err := json.Marshal(map[string]string{"token": claimToken.AccessToken})
	if err != nil {
		fmt.Println("Error marshalling payload:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Error creating request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
	// Client := &http.Client{}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
	if resp.StatusCode != 200 {
		fmt.Println("Invalid token")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid token"))
		return
	}
	defer resp.Body.Close()

	username, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	var user models.User
	c.DB.Where(models.User{Username: string(username)}).FirstOrCreate(&user)
	fmt.Println(user)

	tokenSession, expiredTime, err := handleJWT(user.ID)
	if err != nil {
		fmt.Println("Error creating session token:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	err = CachingUser(user.ID, tokenSession, c.RDB, expiredTime)
	if err != nil {
		fmt.Println("Error storing token to cache:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenSession)
}

func (c controller) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	tokenString := r.FormValue("session")
	if tokenString == "" {
		fmt.Println("No token provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No token provided"))
		return
	}
	Md5TokenString := fmt.Sprintf("%x", md5.Sum([]byte(tokenString)))
	if err := c.RDB.Del(context.Background(), "session:"+Md5TokenString).Err(); err != nil {
		fmt.Println("Error deleting token from cache:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logged out successfully")
}

func GetOauthURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(os.Getenv("OAUTH_URL"))
}

func GetOauthClientURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(os.Getenv("OAUTH_CLIENT_URL"))
}
