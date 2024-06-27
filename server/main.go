package main

import (
	"log"
	"net/http"

	"github.com/bigbluewhale111/rest_api/cache"
	"github.com/bigbluewhale111/rest_api/controllers"
	"github.com/bigbluewhale111/rest_api/db"
	"github.com/gorilla/mux"
)

func main() {
	DB := db.Init()
	RDB := cache.Init()
	c := controllers.New(DB, RDB)
	router := mux.NewRouter()

	unauthenticatedSubRouter := router.PathPrefix("/auth").Subrouter()
	unauthenticatedSubRouter.HandleFunc("/callback", c.Callback).Methods("GET")
	unauthenticatedSubRouter.HandleFunc("/getOauthURL", controllers.GetOauthURL).Methods("GET")

	authenticatedSubRouter := router.PathPrefix("/api").Subrouter()
	authenticatedSubRouter.Use(c.Authenticate)
	authenticatedSubRouter.HandleFunc("/tasks", c.GetAllTasks).Methods("GET")
	authenticatedSubRouter.HandleFunc("/create/task", c.AddTask).Methods("POST", "OPTIONS")
	authenticatedSubRouter.HandleFunc("/task/{id}", c.GetTask).Methods("GET")
	authenticatedSubRouter.HandleFunc("/edit/task/{id}", c.UpdateTask).Methods("POST", "OPTIONS")
	authenticatedSubRouter.HandleFunc("/delete/task/{id}", c.DeleteTask).Methods("GET")
	authenticatedSubRouter.HandleFunc("/logout", c.Logout).Methods("GET")

	log.Println("API is running on port 7001...")
	http.ListenAndServe(":7001", router)
}
