package main

import (
	"log"
	"net/http"

	"github.com/bigbluewhale111/rest_api/controllers"
	"github.com/bigbluewhale111/rest_api/db"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	DB := db.Init()
	c := controllers.New(DB)
	router.HandleFunc("/task_api/tasks", c.GetAllTasks).Methods("GET")
	router.HandleFunc("/task_api/create/task", c.AddTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/task_api/task/{id}", c.GetTask).Methods("GET")
	router.HandleFunc("/task_api/edit/task/{id}", c.UpdateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/task_api/delete/task/{id}", c.DeleteTask).Methods("GET")
	log.Println("API is running on port 3000...")
	http.ListenAndServe(":3000", router)
}
