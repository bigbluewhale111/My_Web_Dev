package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bigbluewhale111/rest_api/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var client_id = os.Getenv("CLIENT_ID")         // CREDS
var client_secret = os.Getenv("CLIENT_SECRET") // CREDS
var secret = os.Getenv("SECRET")               // CREDS

// func addCorsHeader(res http.ResponseWriter) {
// 	headers := res.Header()
// 	headers.Add("Access-Control-Allow-Origin", "*")
// 	headers.Add("Vary", "Origin")
// 	headers.Add("Vary", "Access-Control-Request-Method")
// 	headers.Add("Vary", "Access-Control-Request-Headers")
// 	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
// 	headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
// }

func (c controller) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authorId := r.Context().Value(UserKey("user")).(models.User).ID

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

	authorId := r.Context().Value(UserKey("user")).(models.User).ID

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

	authorId := r.Context().Value(UserKey("user")).(models.User).ID

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

	authorId := r.Context().Value(UserKey("user")).(models.User).ID

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

	authorId := r.Context().Value(UserKey("user")).(models.User).ID

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

func (c controller) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// get code params from callback
	params := r.URL.Query()
	code := params.Get("code")
	fmt.Println(code)

	// get access token
	jsonData := []byte(`{"client_id":"` + client_id + `","client_secret":"` + client_secret + `","code":"` + code + `"}`)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	Client := &http.Client{Transport: tr}
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	var oauthResponse models.OauthResponse
	json.Unmarshal(body, &oauthResponse)
	fmt.Println(oauthResponse)

	// get User info
	url := "https://api.github.com/applications/" + client_id + "/token"
	payload, err := json.Marshal(map[string]string{"access_token": oauthResponse.AccessToken})
	if err != nil {
		fmt.Println("Error marshalling payload:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	rreq, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Error creating request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	rreq.Header.Set("Accept", "application/vnd.github+json")
	rreq.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	rreq.SetBasicAuth(client_id, client_secret)
	resp, err = Client.Do(rreq)
	if err != nil {
		fmt.Println("Error sending request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	var tokenResp models.TokenResponse
	json.Unmarshal(body, &tokenResp)
	fmt.Println(tokenResp.Auth_User)
	var githubUser models.GithubUser = tokenResp.Auth_User

	// save user info to db
	var user models.User
	c.DB.Where(models.User{GithubId: githubUser.Id}).Assign(models.User{Username: githubUser.Login, AccessToken: oauthResponse.AccessToken}).FirstOrCreate(&user)
	fmt.Println(user)
	// create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTClaim{
		Username: user.Username,
		Id:       user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error creating token:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	fmt.Println(tokenString)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenString)
}

func (c controller) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	url := "https://api.github.com/applications/" + client_id + "/token"

	userAcessToken := r.Context().Value(UserKey("user")).(models.User).AccessToken

	payload, err := json.Marshal(map[string]string{"access_token": userAcessToken})
	if err != nil {
		fmt.Println("Error marshalling payload:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Error creating request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.SetBasicAuth(client_id, client_secret)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	Client := &http.Client{Transport: tr}
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	if resp.StatusCode != 204 {
		fmt.Println("Error deleting token:" + strconv.Itoa(resp.StatusCode))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	// clear cookie
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logged out successfully")
}

func (c controller) LoginAs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "")
}

func GetClientID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client_id)
}
