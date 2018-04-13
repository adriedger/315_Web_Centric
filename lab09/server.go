// CMPT 315 (Winter 2018)
// Lab 09: JSON Web Tokens
// Author: Nicholas M. Boers
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	// use "git pull" to obtain the database driver; compiling requires it, but you'll only be interacting with the
	// database through part 3
	_ "github.com/mattn/go-sqlite3"
)

// provideSecrets is a handler that should only be accessible to authenticated clients, i.e., those clients that possess
// a valid JWT
//
// Note: you should *not* need to modify this function.
func provideSecrets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "the secret is: 42\n")
}

// createRoutes creates both authenticated and unauthenticated routes for the web service
//
// Note: you should *not* need to modify this function.
func createRoutes(db *sql.DB) *mux.Router {
	authenticator := NewAuthenticator()

	router := mux.NewRouter()

	// create route for obtaining a JWT
	router.Handle("/api/v1/tokens", authenticator).Methods("POST")

	// create routes that require a valid JWT
	authRouter := router.PathPrefix("/api/v1/").Subrouter()
	authRouter.Use(authenticator.Middleware)
	authRouter.HandleFunc("/secrets", provideSecrets).Methods("GET")

	// add a route for serving files
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("htdocs/"))))

	return router
}

// openDatabase opens a database for this lab; it is only necessary for the optional part 3
//
// Note: you should *not* need to modify this function.
func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "lab09.sqlite3")
	if err != nil {
		return nil, err
	}

	// causes the creation of the sqlite3 file
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := openDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: cannot open database: %v\n", err)
		os.Exit(1)
	}

	router := createRoutes(db)
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}
