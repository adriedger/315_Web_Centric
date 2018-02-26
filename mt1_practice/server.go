//look at time package

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var db *Database
var logger *log.Logger

func init() {
	var quiet bool
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&quiet, "q", false, "hide log messages")
	flag.Parse()
	if quiet {
		logger = log.New(ioutil.Discard, "Status log: ", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "Status log: ", log.LstdFlags)
	}
}

func handleGetTime(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t)
	err := json.NewEncoder(w).Encode(t)
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	db, err = OpenDatabase()
	if err != nil {
		log.Fatalf("OpenDatabase: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", handleGetTime)
	http.ListenAndServe(":8080", nil)
	logger.Printf("Success")
}
