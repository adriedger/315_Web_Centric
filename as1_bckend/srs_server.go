package main

import (
	"encoding/json"
	//    "encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	//	"sync"
	"time"
)

var db *Database

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randomStringGen(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

type Class struct {
	ClassID    string `db:"class_id"`
	ClassName  string `db:"class_name"`
	CreatorKey string `db:"creator_key"`
}

type Enrollment struct {
	EnrollID string `db:"enroll_id"`
	Username string `db:"username"`
	ClassID  string `db:"class_id"`
}

type Question struct {
	Question   string `db:"question"`
	ClassID    string `db:"class_id"`
	Answer     string `db:"answer"`
	KeyAttempt string `db:"key_attempt`
}

type Response struct {
	Answer   string `db:"response"`
	Question string `db:"question"`
	EnrollID string `db:"enroll_id"`
}

func handleGetHome(w http.ResponseWriter, r *http.Request) {
	home := struct {
		CreateClassLink string
		JoinClassLink   string
	}{"/api/v1/classes/create", "/api/v1/classes/join"}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(home)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	class.ClassID = randomStringGen(4)
	class.CreatorKey = randomStringGen(4)

	// loop until no same id error
	err = db.AddClass(class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetClass(w http.ResponseWriter, r *http.Request) {
	mURLVars := mux.Vars(r)
	//returs a Class struct
	class, err := db.GetClass(mURLVars["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		fmt.Println("error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleJoinClass(w http.ResponseWriter, r *http.Request) {
	mURLVars := mux.Vars(r)
	var enrollment Enrollment
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	enrollment.EnrollID = randomStringGen(4)
	enrollment.ClassID = mURLVars["id"]
	//fmt.Printf("%+v", enrollment)
	err = db.JoinClass(enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	//fmt.Println("JOINED")
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreateQuestion(w http.ResponseWriter, r *http.Request) {
	mURLVars := mux.Vars(r)
	var question Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	question.ClassID = mURLVars["id"]
	//fmt.Printf("%+v", question)
	err = db.AddQuestion(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("QUESTION CREATED")
}

func handleAddResponse(w http.ResponseWriter, r *http.Request) {
	var response Response
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	//fmt.Printf("%+v", response)
	err = db.AddResponse(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("RESPONSE ADDED")
}

func handleModifyResponse(w http.ResponseWriter, r *http.Request) {
	var response Response
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	//fmt.Printf("%+v", response)
	err = db.ModifyResponse(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("RESPONSE MODIFIED")
}

func handleDeleteQuestion(w http.ResponseWriter, r *http.Request) {
	//gotta delete all responses assosiated with question
	//need key attemp to compare with class creator key, classid and question
	var question Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%+v", question)
	err = db.DeleteQuestion(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("QUESTION DELETED")
}

func main() {
	var err error
	db, err = OpenDatabase()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer db.Close()
	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()
	router.HandleFunc("/api/v1", handleGetHome).Methods("GET")
	router.HandleFunc("/api/v1/classes/create", handleCreateClass).Methods("POST")
	router.HandleFunc("/api/v1/classes/{id:(?:[0-9]|[A-Z]){4}}", handleGetClass).Methods("GET")
	router.HandleFunc("/api/v1/classes/join/{id:(?:[0-9]|[A-Z]){4}}", handleJoinClass).Methods("POST")
	router.HandleFunc("/api/v1/questions/create/{id:(?:[0-9]|[A-Z]){4}}", handleCreateQuestion).Methods("POST")
	router.HandleFunc("/api/v1/responses/add", handleAddResponse).Methods("POST")
	router.HandleFunc("/api/v1/responses/modify", handleModifyResponse).Methods("POST")
	router.HandleFunc("/api/v1/questions/delete", handleDeleteQuestion).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
