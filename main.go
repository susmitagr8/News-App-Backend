package main

import (
	"log"
	"net/http"

	"news-app-backend/store"

	"github.com/gorilla/handlers"
)

func main() {

	router := store.NewRouter() // create routes
	// These two lines are important if you're designing a front-end to utilise this API methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	// Launch server with CORS validations
	log.Fatal(http.ListenAndServe(":"+"8081", handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
