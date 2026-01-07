package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//define a home handler function which writes a byte slice
//containing "hello from snippetbox" as the response body.

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetbox"))

}

// add a snippetview handler function
func snippetView(w http.ResponseWriter, r *http.Request) {

	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the fmt.Sprintf() function to interpolate the id value with a
	// message, then write it as the HTTP response.
	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))
}

// add a snippetcreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a form for creating a new snippet..."))
}
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Save a new snippet..."))
}

func main() {
	//use the http.NewServeMux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)

	//Restricting subtree paths
	// 	So if you have the route pattern "/{$}", it effectively means match a single slash, followed
	// by nothing else. It will only match requests where the URL path is exactly /.

	// what happen is the system match only the slash'/' nothing else
	//so if a slash meets any unknown route for ex '/foo'
	//it will redirect it to '/' endpoint thats why we made the above.

	// 	Request URL paths are automatically sanitized. If the request path contains any . or ..
	// elements or repeated slashes, the user will automatically be redirected to an equivalent
	// clean URL. For example, if a user makes a request to /foo/bar/..//baz they will
	// automatically be sent a 301 Permanent Redirect to /foo/baz instead.

	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	//add the {id} wildcard segment
	//{id} must match a non-empty path segment
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
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

// Important: Before we continue, I should explain that Go’s servemux treats the route
// pattern "/" like a catch-all. So at the moment all HTTP requests to our server will be
// handled by the home function, regardless of their URL path. For instance, you can visit
// a different URL path like http://localhost:4000/foo/bar and you’ll receive exactly
// the same response
