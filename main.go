package main

import (
	"log"
	"net/http"
)

//define a home handler function which writes a byte slice
//containing "hello from snippetbox" as the response body.

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetbox"))

}

func main() {
	//use the http.NewServeMux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	//print a log message to say that the server is starting
	log.Print("starting server on :4000")

	// use the http.ListenAndServe() function to start a new web server.
	//we pass in two parameters:
	//the tcp network address to listen on(in this case ":4000")
	//and the servemux we just created. if http.ListenAndServe() returns
	//an error we use the log.Fatal() function to log the error message
	// that any error returned by http.ListenAndServe() is always non-nil.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// Note: The home handler function is just a regular Go function with two parameters.
// The http.ResponseWriter parameter provides methods for assembling a HTTP
// response and sending it to the user, and the *http.Request parameter is a pointer to
// a struct which holds information about the current request (like the HTTP method
// and the URL being requested). We’ll talk more about these parameters and
// demonstrate how to use them as we progress through the book.
