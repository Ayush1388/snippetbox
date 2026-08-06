package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("save a new snippet"))
}

func main() {
	//use the http.NewServeMux() function to initialize a new servemux,
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/createPost", snippetCreatePost)

	log.Print("starting server on :4000")

	//use the http.ListenAndServe() function to start a new web server
	//pass two parameters the tcp network addr to listen on
	//(":4000") and the servemux we just created.
	//if returns an error we will use log.fatal() function to log the error
	//message and exit
	//the error returns is always non-nil.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
