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

type Class struct {
	Name       string
	ID         string
	CreatorKey string
	Joined     []string
	//	Questions	[]Question
}

func handleJoinClass(w http.ResponseWriter, r *http.Request) {
	type Entry struct {
		ClassID  string
		Username string
	}
	var entry Entry
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%+v", entry)
	// if entry.ClassID == ID in database, gen id for user, add user to class arr of users
	// return status of join, class name and id
}

func handleGetClass(w http.ResponseWriter, r *http.Request) {
	mURLVars := mux.Vars(r)

	mutex.Lock()
	if mURLVars["id"] == "9001" {
		fmt.Println("over 9000!")
	} else {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	}
	mutex.Unlock()
}

func handleCreateClass(w http.ResponseWriter, r *http.Request) {
	var class Class
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	class.ID = randomStringGen(4)
	class.CreatorKey = randomStringGen(4)
	//fmt.Printf("%+v", classes)
	js, err := json.Marshal(class)
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
	rand.Seed(time.Now().UnixNano())
	router := mux.NewRouter()
	router.HandleFunc("/api/v1", handleGetHome).Methods("GET")
	router.HandleFunc("/api/v1/classes/create", handleCreateClass).Methods("POST")
	router.HandleFunc("/api/v1/classes/{id:(?:[0-9]|[A-Z]){4}}", handleGetClass).Methods("GET")
	router.HandleFunc("/api/v1/classes/join", handleJoinClass).Methods("POST")

	http.ListenAndServe(":8080", router)
}
