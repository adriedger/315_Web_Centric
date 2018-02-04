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

type Student struct {
	StudentName string `db:"student_name"`
}

type Class struct {
	ClassID    string `db:"class_id"`
	ClassName  string `db:"class_name"`
	CreatorKey string `db:"creator_key"`
}

type Enrollment struct {
	EnrollID    int    `db:"enrollment"`
	ClassID     string `db:"class_id"`
	StudentName string `db:"student_name"`
}

type Question struct {
	Text    string `db:"question"`
	ClassID string `db:"class_id"`
	Answer  string `db:"answer"`
}

//responces to choose from for each question
type Response struct {
	Text        string `db:"responce"`
	Question    string `db:"question"`
	StudentName string `"db:"student_name"`
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
	var student Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	//fmt.Printf("%+v", student)
	//run db
	err = db.JoinClass(mURLVars["id"], student.StudentName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("JOINED")
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

	http.ListenAndServe(":8080", router)
}
