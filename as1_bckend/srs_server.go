package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randomStringGen(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

type User struct {
	Name string
	ID   string
}

type Class struct {
	Name    string
	ID      string
	Creator User
	//	Joined  []User
	//	Questions	[]Question

}

func handlePostClass(w http.ResponseWriter, r *http.Request) {
	var jsonBlob = []byte(`[{"Name": "CMPT 101"}]`)
	var classes []Class
	err := json.Unmarshal(jsonBlob, &classes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	classes[0].ID = randomStringGen(4)
	classes[0].Creator = User{"yo", randomStringGen(4)}
	//fmt.Printf("%+v", classes)
	js, err := json.Marshal(classes[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleGetHome(w http.ResponseWriter, r *http.Request) {
	type HomeOptions struct {
		CreateClassLink string
		JoinClassLink   string
	}
	home := HomeOptions{"/api/v1/classes?q=create", "/api/v1/classes?q=join"}
	js, err := json.Marshal(home)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", handleGetHome).Methods("GET")
	router.HandleFunc("/api/v1/classes/create/", handlePostClass).Methods("POST")

	http.ListenAndServe(":8080", router)
}
