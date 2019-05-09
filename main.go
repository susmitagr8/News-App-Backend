package main

import (
	"net/http"

	"news-app/store"

	"github.com/rs/cors"
)

func main() {
	router := store.NewRouter()
	// router.HandleFunc("/xx", xyz).Methods("OPTIONS")
	// router := store.NewRouter() // create routes
	// // These two lines are important if you're designing a front-end to utilise this API methods
	// allowedHeaders := handlers.AllowedHeaders([]string{"*"})
	// allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	// allowedMethods := handlers.AllowedMethods([]string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut})
	// allowedCre := handlers.AllowCredentials()
	// // Launch server with CORS validations
	// log.Fatal(http.ListenAndServe(":"+"8081", handlers.CORS(allowedHeaders, allowedCre, allowedOrigins, allowedMethods)(router)))

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})

	http.ListenAndServe(":8081", corsOpts.Handler(router))
}
