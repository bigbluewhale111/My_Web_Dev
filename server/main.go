package main

import (
	"log"
	"net/http"

	"github.com/bigbluewhale111/rest_api/controllers"
	"github.com/bigbluewhale111/rest_api/db"
	"github.com/bigbluewhale111/rest_api/env"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	env.LoadEnvVars()
	DB := db.Init()
	c := controllers.New(DB)
	router.HandleFunc("/api/tasks", c.GetAllTasks).Methods("GET")
	router.HandleFunc("/api/create/task", c.AddTask).Methods("POST")
	router.HandleFunc("/api/task/{id}", c.GetTask).Methods("GET")
	router.HandleFunc("/api/task/{id}", c.UpdateTask).Methods("PUT")
	router.HandleFunc("/api/task/{id}", c.DeleteTask).Methods("DELETE")
	log.Println("API is running on port 3000...")
	http.ListenAndServe(":3000", router)
}
