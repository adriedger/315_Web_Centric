package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
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
	ClassName  string `db:"class_name"`
	CreatorKey string `db:"creator_key"`
}

type Enrollment struct {
	Username  string `db:"username"`
	ClassName string `db:"class_name"`
}

type Question struct {
	Question   string `db:"question"`
	Answer     string `db:"answer"`
	ClassName  string `db:"class_name"`
	KeyAttempt string `db:"key_attempt`
}

type Response struct {
	Response  string `db:"response"`
	Question  string `db:"question"`
	ClassName string `db:"class_name"`
	Username  string `db:"username"`
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
	class.CreatorKey = randomStringGen(4)

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

func handleJoinClass(w http.ResponseWriter, r *http.Request) {
	var enrollment Enrollment
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	//fmt.Printf("%+v", enrollment)
	err = db.JoinClass(enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddQuestion(w http.ResponseWriter, r *http.Request) {
	var question Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
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
	err = db.ModifyResponse(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("RESPONSE MODIFIED")
}

func handleDeleteQuestion(w http.ResponseWriter, r *http.Request) {
	var question Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	//fmt.Printf("%+v", question)
	err = db.DeleteQuestion(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	fmt.Println("QUESTION DELETED")
}

func handleGetQuestions(w http.ResponseWriter, r *http.Request) {
	mURLVars := mux.Vars(r)
	questions, err := db.GetQuestions(mURLVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(questions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetResponses(w http.ResponseWriter, r *http.Request) {
	//get keyattempt
	var question Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	responses, err := db.GetResponses(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	router.HandleFunc("/api/v1/classes/join", handleJoinClass).Methods("POST")
	router.HandleFunc("/api/v1/questions/create", handleAddQuestion).Methods("POST")
	router.HandleFunc("/api/v1/responses/add", handleAddResponse).Methods("POST")
	router.HandleFunc("/api/v1/responses/modify", handleModifyResponse).Methods("POST")
	router.HandleFunc("/api/v1/questions/delete", handleDeleteQuestion).Methods("DELETE")
	router.HandleFunc("/api/v1/classes/questions/{name}", handleGetQuestions).Methods("GET")
	router.HandleFunc("/api/v1/questions/responses", handleGetResponses).Methods("GET")

	http.ListenAndServe(":8080", router)
}
