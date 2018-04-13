// CMPT 315 (Winter 2018)
// Lab 09: JSON Web Tokens
// Author: Nicholas M. Boers
package main

import (
	"fmt"
	"net/http"
)

type Authenticator struct {
	// one could add additional fields here, e.g., a pointer to a database
}

// NewAuthenticator creates a new authenticator; the authenticator provides a ServeHTTP method for obtaining tokens and
// a Middleware method for verifying tokens
func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

// ServeHTTP verifies supplied credentials; if the credentials are valid, it returns a JWT in the body; otherwise, it
// returns 400 (Bad Request)
func (a *Authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("authenticate handler: unimplemented")
}

// Middleware returns a handler that verifies the JWT included in the "Authorization" header; if the JWT is valid, it
// calls the next handler; otherwise, it returns 401 (Unauthorized) along with a "WWW-Authenticate" header indicating
// the needed for a bearer token
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("authenticate middleware: unimplemented")
		next.ServeHTTP(w, r)
		return
	})
}
