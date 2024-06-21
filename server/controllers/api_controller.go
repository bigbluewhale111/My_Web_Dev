package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bigbluewhale111/rest_api/models"
	"github.com/gorilla/mux"
)

var testAuthorId uint32 = 0

func setAccessControlAllowOrigin(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (c controller) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	setAccessControlAllowOrigin(w)
	var tasks []models.Task
	if result := c.DB.Find(&tasks, "author_id = ?", testAuthorId); result.Error != nil {
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
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	setAccessControlAllowOrigin(w)
	if err != nil {
		fmt.Println("Error reading body:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	var newtask models.NewTask
	json.Unmarshal(body, &newtask)

	if result := c.DB.Create(&models.Task{
		Name:        newtask.Name,
		Description: newtask.Description,
		Status:      newtask.Status,
		DueDate:     newtask.DueDate,
		AuthorID:    testAuthorId,
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
	vars := mux.Vars(r)
	var id, _ = strconv.Atoi(vars["id"])
	setAccessControlAllowOrigin(w)
	var task models.Task
	if result := c.DB.First(&task, "id = ? AND author_id = ?", id, testAuthorId); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (c controller) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id, _ = strconv.Atoi(vars["id"])

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	setAccessControlAllowOrigin(w)
	if err != nil {
		fmt.Println("Error reading body:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	var UpdateTask models.NewTask
	json.Unmarshal(body, &UpdateTask)

	var task models.Task
	if result := c.DB.First(&task, "id = ? AND author_id = ?", id, testAuthorId); result.Error != nil {
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
	if UpdateTask.Status != "" {
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
	vars := mux.Vars(r)
	var id, _ = strconv.Atoi(vars["id"])
	setAccessControlAllowOrigin(w)
	var task models.Task
	if result := c.DB.First(&task, "id = ? AND author_id = ?", id, testAuthorId); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	if result := c.DB.Delete(&task, "id = ? AND author_id = ?", id, testAuthorId); result.Error != nil {
		fmt.Println("Error deleting task:" + result.Error.Error())
		setAccessControlAllowOrigin(w)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Add("Content-type", "application/json")
	setAccessControlAllowOrigin(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Task deleted successfully")
}
