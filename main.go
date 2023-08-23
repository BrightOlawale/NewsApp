package main

import (
	"fmt"
	"github.com/BrightOlawale/NewsApp/news"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Go provides two template packages in its standard library: text/template and html/template
// The html/template package is used to generate output that is safe against code injection.
// tpl is a package level variable that points to a template definition from the provided files.
// The call to template.ParseFiles parses the index.html file in the root of our project directory
// and validates it. The invocation of template.ParseFiles is wrapped with template.Must so that
// the code panics if an error is obtained while parsing the template file.
var tpl = template.Must(template.ParseFiles("index.html"))

// Search : Represents each query made by the user
type Search struct {
	Query     string
	NextPage  int
	TotalPage int
	Results   *news.Results
}

// indexHandler : serves as an index handler function for an HTTP request
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// NB: The w parameter is the structure we use to send responses to an HTTP request.
// It implements a Write() method which accepts a slice of bytes and writes
// the data to the connection as part of an HTTP response.

// The r parameter represents the HTTP request received from the client.
// It’s how we access the data sent by a client to the server.

// searchHandler : serves as a search handler function for an HTTP request.
// It takes the newsAPI argument and returns a http handler.
func searchHandler(newsAPI *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())

		// Check if there was an error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// The Query() method of the url.URL type returns a map of the query string
		// parameters. The map keys are the query string parameter names and the
		// map values are the corresponding query string parameter values.
		queryParams := u.Query()

		// The Get() method of the url.Values type returns the first value associated
		// with the given key. If there are no values associated with the key, the
		// Get() method returns an empty string.
		// q represents the user’s query, and page is used to page through the results
		searchQuery := queryParams.Get("q")
		page := queryParams.Get("page")

		// Check if page number was retrieved and if not set it to 1
		if page == "" {
			page = "1"
		}

		results, err := newsAPI.FetchEverything(searchQuery, page)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v", results)

		//nextPage, err := strconv.Atoi(page)
		//
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		//search := &Search{
		//	Query:     searchQuery,
		//	NextPage:  nextPage,
		//	TotalPage: int(math.Ceil(float64(results.TotalResults) / float64(newsAPI.PageSize))),
		//	Results:   results,
		//}
	}
}

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

	// Retrieve the API key from the environment
	apiKey := os.Getenv("NEWS_API_KEY")

	// Check if API key was retrieved
	if apiKey == "" {
		log.Fatal("API key not found in .env file")
	}

	// The http.Client type is used to make HTTP requests. It provides a default
	// configuration that is suitable for most use cases. You can create a new
	// http.Client instance by calling the http.Client{} constructor.
	// The Timeout field of the http.Client type is used to set the maximum
	// amount of time a request can take before it is canceled.
	myClient := &http.Client{Timeout: 10 * time.Second}

	// The NewClient() function is used to create a new client instance used to
	// make requests to the News API. It accepts a http.Client instance, an API
	// key, and a page size as parameters and returns a pointer to a Client instance.
	newsAPI := news.NewClient(myClient, apiKey, 20)

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
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// You then use the mux's HandleFunc method to associate specific paths
	// with their respective handler functions
	mux.HandleFunc("/", indexHandler)

	// Register the searchHandler function as the handler function for the /search path
	mux.HandleFunc("/search", searchHandler(newsAPI))

	// "http.ListenAndServe" function is used to start an HTTP server that listens for
	// incoming requests and serves them using a specified handler
	http.ListenAndServe(":"+port, mux)
}
