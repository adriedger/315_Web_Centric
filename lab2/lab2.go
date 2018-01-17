// CMPT 315 (Winter 2018)
// Lab #2: Introduction to gorilla/mux
//
// For this lab, use gorilla/mux to implement the following two routes:
//
// (1) GET requests to /api/v1/digits/n cause a lookup in the digits map, and if an entry
//     exists, cause the server to print the the corresponding string to the ResponseWriter.
//     On success: 200 (OK)
//     On failure: 410 (Gone)
//
// (2) DELETE requests to /api/v1/digits/n cause a lookup in the digits map, and if an entry
//     exists, cause the server to delete the entry from the map.
//     On success: 204 (No Content)
//     On failure: 410 (Gone)
//
// All other routes may generate a 404 (Not Found).

package main

import (
    "fmt"
	"net/http"
	"sync"
    "github.com/gorilla/mux"
    "io"
)

// mutex exists to ensure safe access to the map "digits"
var mutex sync.Mutex

// digits maps from integers to their corresponding strings
var digits map[int]string = map[int]string{
	0: "zero",
	1: "one",
	2: "two",
	3: "three",
	4: "four",
	5: "five",
	6: "six",
	7: "seven",
	8: "eight",
	9: "nine",
}

// handleDeleteDigit handles DELETE requests to /api/v1/digits/n by removing the digit
// from the map (if it exists). If the digit exists (and is removed), it produces a header
// containing 204 (No Content). If the digit doesn't exist, it produces a header
// containing 410 (Gone).
func handleDeleteDigit(w http.ResponseWriter, r *http.Request) {

	// code goes here
}

// handleGetDigit handles GET requests to /api/v1/digits/n by attempting to get the
// digit's string from the map. If the digit exists, it gets the string and writes that
// string in the response body, which will produce a header containing 200 (OK). If the
// digit doesn't exist, it produces a header containing 410 (Gone).
func handleGetDigit(w http.ResponseWriter, r *http.Request) {

    // code goes here
}

// main creates a gorilla/mux router that handles the two routes described in the
// introduction
func main() {
    r := mux.NewRouter()
    r.HandleFunc("/api/v1/digits/",handleDeleteDigit).Methods("GET")
    err := http.ListenAndServe(":8080", r)
    if err != nil{
        fmt.Println(err)
    }
	// code goes here
}
