package main

import (
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Go provides two template packages in its standard library: text/template and html/template
// The html/template package is used to generate output that is safe against code injection.
// tpl is a package level variable that points to a template definition from the provided files.
// The call to template.ParseFiles parses the index.html file in the root of our project directory
// and validates it. The invocation of template.ParseFiles is wrapped with template.Must so that
// the code panics if an error is obtained while parsing the template file.
var tpl = template.Must(template.ParseFiles("index.html"))

// indexHandler : serves as an index handler function for an HTTP request
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// NB: The w parameter is the structure we use to send responses to an HTTP request.
// It implements a Write() method which accepts a slice of bytes and writes
// the data to the connection as part of an HTTP response.

// The r parameter represents the HTTP request received from the client.
// It’s how we access the data sent by a client to the server.

func main() {
	// The Load method reads the .env file and loads the set variables into
	// the environment so that they can be accessed through the os.Getenv() method.
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")

	// Check if port number was retrieved
	if port == "" {
		port = "8080"
	}

	// The below line instantiates a file server object by passing the directory
	// where all our static files are placed
	fs := http.FileServer(http.Dir("assets"))

	// A multiplexer (mux) is an HTTP request multiplexer used to route
	// incoming HTTP requests to their corresponding handler functions.
	// A ServeMux is a central router that helps map incoming
	// HTTP requests to their corresponding handler functions
	mux := http.NewServeMux()

	// StripPrefix() will cut the /assets/ part and forward the modified
	// request to the handler returned by http.FileServer() so it will
	// see the requested resource as style.css
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	// You then use the mux's HandleFunc method to associate specific paths
	// with their respective handler functions
	mux.HandleFunc("/", indexHandler)

	// "http.ListenAndServe" function is used to start an HTTP server that listens for
	// incoming requests and serves them using a specified handler
	http.ListenAndServe(":"+port, mux)
}
