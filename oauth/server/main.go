package main

import (
	"fmt"
	"net/http"

	"github.com/bigbluewhale111/oauth_server/cache"
	"github.com/bigbluewhale111/oauth_server/controllers"
	"github.com/bigbluewhale111/oauth_server/db"
	"github.com/gorilla/mux"
)

func main() {
	DB := db.Init()
	RDB := cache.Init()
	c := controllers.New(DB, RDB)
	router := mux.NewRouter()
	oauth_client_subrouter := router.PathPrefix("/oauth").Subrouter()
	oauth_client_subrouter.Use(controllers.BasicAuth)
	oauth_client_subrouter.HandleFunc("/token", c.SendToken).Methods("POST")
	oauth_client_subrouter.HandleFunc("/revoke", c.RevokeToken).Methods("POST")

	router.HandleFunc("/login", c.LoginUser).Methods("POST")
	router.HandleFunc("/register", c.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.ServeFileLogin).Methods("GET")
	router.HandleFunc("/register", controllers.ServeFileRegister).Methods("GET")
	router.HandleFunc("/validate", c.ValidateToken).Methods("POST")

	fmt.Println("Server is running on port 7000")
	http.ListenAndServe(":7000", router)
}
