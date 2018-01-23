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
	Joined  []User
	//	Questions	[]Question
}

func handleGetClass(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("yo")
	/*var ids map[string]string = map[string]string {
		"9001": "Alice"
		"9002": "Bob"
		"9003"; "Claire"
	}*/
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
	var jsonBlob = []byte(`[{"Name": "CMPT 101", "Creator": {"Name": "Bob"}}]`)
	var classes []Class
	err := json.Unmarshal(jsonBlob, &classes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	classes[0].ID = randomStringGen(4)
	classes[0].Creator.ID = randomStringGen(4)
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

	http.ListenAndServe(":8080", router)
}
