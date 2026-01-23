package main

import (
	"database/sql" // [NEW] Needed to work with SQL databases
	"flag"
	"log/slog"
	"net/http"
	"os"

	// [NEW] "Why this import has an underscore?"
	// The underscore (_) is an alias. We don't use the 'mysql' package directly in our code,
	// but we need its 'init()' function to run so it registers itself with Go's 'database/sql' package.
	// Without this, sql.Open() wouldn't know how to talk to MySQL.
	"github.com/Ayush1388/Snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {

	//dsn := "web:pass@tcp(127.0.0.1:3306)/snippetbox?parseTime=true"

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// [NEW] Define a new command-line flag for the MySQL DSN string.
	// "Why do we need this?"
	// Ideally, we don't want hardcoded passwords in code. Using a flag allows us to pass
	// the password when running the app (e.g., go run . -dsn="...") or use a default.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// [NEW] Initialize the Database Connection Pool
	// "How does this work?"
	// We call our helper function openDB() (defined at the bottom).
	// It doesn't just "open" a file; it creates a "pool" of connections that can be reused.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// [NEW] Defer closing the database
	// "Why defer?"
	// We want the database connection to stay alive while the application runs.
	// 'defer' ensures that db.Close() is called ONLY when the main() function exits (shuts down).
	defer db.Close()

	// Initialize a models.SnippetModel instance containing the connection pool
	// and add it to the application dependencies.
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	//use the http.NewServeMux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("C:/Users/chand/Desktop/snippetbox/ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)

	//Restricting subtree paths
	//  So if you have the route pattern "/{$}", it effectively means match a single slash, followed
	// by nothing else. It will only match requests where the URL path is exactly /.

	// what happen is the system match only the slash'/' nothing else
	//so if a slash meets any unknown route for ex '/foo'
	//it will redirect it to '/' endpoint thats why we made the above.

	//  Request URL paths are automatically sanitized. If the request path contains any . or ..
	// elements or repeated slashes, the user will automatically be redirected to an equivalent
	// clean URL. For example, if a user makes a request to /foo/bar/..//baz they will
	// automatically be sent a 301 Permanent Redirect to /foo/baz instead.

	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	//add the {id} wildcard segment
	//{id} must match a non-empty path segment
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it (i.e. prefix it with
	// the * symbol) before using it. Note that we're using the log.Printf()
	// function to interpolate the address with the log message.
	//log.Printf("starting server on %s", *addr)

	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", "addr", *addr)
	//print a log message to say that the server is starting
	//log.Print("starting server on :4000")

	// use the http.ListenAndServe() function to start a new web server.
	//we pass in two parameters:
	//the tcp network address to listen on(in this case ":4000")
	//and the servemux we just created. if http.ListenAndServe() returns
	//an error we use the log.Fatal() function to log the error message
	// that any error returned by http.ListenAndServe() is always non-nil.

	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	err = http.ListenAndServe(*addr, mux)

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
	//err := http.ListenAndServe(":4000", mux)
	// log.Fatal(err)
}

// [NEW] The Helper Function
// "Why separate this?"
// It keeps main() clean. This function handles the "dirty work" of:
// 1. Creating the connection pool (sql.Open)
// 2. Verifying the connection actually works (db.Ping)
// If we just used sql.Open(), it wouldn't tell us if the password was wrong until we tried to run a query.
// db.Ping() forces a check immediately.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
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
