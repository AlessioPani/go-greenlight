package main

import (
	"net/http"
)

// healthcheckHandler is the handler that shows information about the application.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Define data map to be converted in JSON.
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	// Convert data in JSON with the specified status code using an helper method and
	// checks for errors.
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
