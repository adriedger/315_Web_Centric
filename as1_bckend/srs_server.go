package main

import (
	"encoding/json"
	//    "encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	//	"sync"
	//	"github.com/jmoiron/sqlx"
	"time"
)

//var mutex sync.Mutex
var db Database

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randomStringGen(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

type Student struct {
	StudentID   int    `db:"student_id"`
	StudentName string `db:"studnet_name"`
}

//classid and creator are random 4 char strings, class id must be unique
type Class struct {
	ClassID    string `db:"class_id"`
	ClassName  string `db:"class_name"`
	CreatorKey string `db:"creator_key"`
	//	Questions	[]Question
}

type Enrollment struct {
	ID        int    `db:"enrollment"`
	StudentID int    `db:"student_id"`
	ClassID   string `db:"class_id"`
}

/*
type Questions struct {
    ID          int
    ClassID     string
    QuestionID  string
}
*/
func handleJoinClass(w http.ResponseWriter, r *http.Request) {
	type Entry struct {
		ClassID     string
		StudentName string
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

	//	mutex.Lock()
	if mURLVars["id"] == "9001" {
		fmt.Println("over 9000!")
	} else {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	}
	//	mutex.Unlock()
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
	/*
		add class to database, loop until no same id error
	*/
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

func main() {
	db, err := OpenDatabase()
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
	router.HandleFunc("/api/v1/classes/join", handleJoinClass).Methods("POST")

	http.ListenAndServe(":8080", router)
}
