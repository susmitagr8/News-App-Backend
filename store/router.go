package store

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var controller = &Controller{repository: Repository{}, MessageRepo: MessageRepo{}}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		"Authentication",
		http.MethodPost,
		"/get-token",
		controller.GetToken,
	},
	Route{
		"AddUser",
		http.MethodPost,
		"/AddUser",
		controller.AddUser,
	},
	Route{
		"AddMessage",
		http.MethodPost,
		"/messageadd",
		AuthenticationMiddleware(controller.AddMessage),
	},
	Route{
		"GetMessages",
		http.MethodGet,
		"/messageget",
		AuthenticationMiddleware(controller.FindMessage),
	},
	Route{
		"XX",
		http.MethodOptions,
		"/messageadd",
		controller.Xyz,
	},
	// Route{
	// 	"GetMessages",
	// 	http.MethodGet,
	// 	"/messageget",
	// 	AuthenticationMiddleware(controller.FindMessage),
	// },
}

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	log.Println("oooooooooooooooo")
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
