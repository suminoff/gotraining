// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to work with the Gorilla Context package.
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/urfave/negroni"
)

// User represents a user in the system.
type User struct {
	Username string
}

// indexHandler handles the index route request.
func indexHandler(res http.ResponseWriter, req *http.Request) {

	// We expect to find a user for this request under
	// the specified key.
	u := context.Get(req, "current_user").(*User)

	// Send the user's username as the response.
	res.Write([]byte(u.Username))
}

// userHandler create the user for the request and saves this user
// inside the context for later use.
func userHandler(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// Create a user value.
	u := User{"mary-jane"}

	// Add the user to the context for this request
	// under the specified key.
	context.Set(req, "current_user", &u)

	// Call the handler that was provided.
	next(res, req)
}

// App loads the API and the middleware.
func App() http.Handler {

	// Create a new mux for handling routes.
	m := http.NewServeMux()

	// Bind the root route to the indexHandler.
	m.HandleFunc("/", indexHandler)

	// Create a new negroni to apply middleware.
	n := negroni.New()

	// Add the middleware handler.
	n.UseFunc(userHandler)

	// Apply the mux to negroni and wrap it around the ClearHandler
	// which clears request values at the end of a request lifetime.
	n.UseHandler(context.ClearHandler(m))

	return n
}

func main() {

	// Start the http server to handle the request for
	// both versions of the API.
	log.Fatal(http.ListenAndServe(":3000", App()))
}
