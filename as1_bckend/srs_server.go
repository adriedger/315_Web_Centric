package main

import (
	"encoding/json"
	//	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

var mutex sync.Mutex

func handleGetHome(w http.ResponseWriter, r *http.Request) {
	type HomeOptions struct {
		CreateClassURL string
		JoinClassURL   string
	}
	home := HomeOptions{"/api/v1/classes/create", "/api/v1/classes/join"}
	js, err := json.Marshal(home)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", handleGetHome).Methods("GET")

	http.ListenAndServe(":8080", router)
}
