package main

import (
	"fmt"
	"net/http"
)

// createMovieHandler is the handler that creates a movie.
// Method: POST
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// showMovieHandler is the handler that shows a movie by its ID.
// Method: GET
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Show details of the movie.
	fmt.Fprintf(w, "show the details of movie %d", id)
}
