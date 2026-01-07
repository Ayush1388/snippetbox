package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

//define a home handler function which writes a byte slice
//containing "hello from snippetbox" as the response body.

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.

	files := []string{
		"C:/Users/chand/Desktop/snippetbox/ui/html/base.html",
		"C:/Users/chand/Desktop/snippetbox/ui/html/partials/nav.html",
		"C:/Users/chand/Desktop/snippetbox/ui/html/pages/home.html",
	}

	// use the template.ParseFiles() function to read the template files into a
	// template set.
	// Notice that we use ... to pass the contents
	// of the files slice as variadic arguments.

	// if there's error , we log the detailed error message, use
	//the http.error() function to send an internatl server error response to the
	//user, and then return from the handler so no subsequent code is executed.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	//then we use the Execute() method on the template set to write the
	//template content as the response body.
	//the last parameter to Execute() represents any dynamic data that we want to pass in ,
	// which for now we'll leave as nil.

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.

	err = ts.ExecuteTemplate(w, "base", nil)
	//err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
	}

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
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// add a snippetcreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a form for creating a new snippet..."))
}
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
