package main

import (
	"fmt"
	"net/http"

	"github.com/bigbluewhale111/oauth_client/controllers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(controllers.HandlingCORS)
	router.HandleFunc("/callback", controllers.CallbackHandler).Methods("GET")
	router.HandleFunc("/logout", controllers.Logout).Methods("GET")

	fmt.Println("Oauth client is running on port 7003")
	http.ListenAndServe(":7003", router)
}
