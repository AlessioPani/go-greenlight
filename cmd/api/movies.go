package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// createMovieHandler is the handler that creates a movie.
// Method: POST
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// showMovieHandler is the handler that shows a movie by its ID.
// Method: GET
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve parameters from URL.
	params := httprouter.ParamsFromContext(r.Context())

	// Get ID.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Show details of the movie.
	fmt.Fprintf(w, "show the details of movie %d", id)
}
