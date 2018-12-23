package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type Profile struct {
	Status string
}

func healthcheck(w http.ResponseWriter, r1 *http.Request) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "newz",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(session)
	profile := Profile{"Healthy"}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Write(js)
}
func main() {
	http.HandleFunc("/healthcheck", healthcheck)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
