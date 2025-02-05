package main

import (
	"net/http"
)

// healthcheckHandler is the handler that shows information about the application.
// Method: GET
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Define data map to be converted in JSON.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	// Convert data in JSON with the specified status code using an helper method and
	// checks for errors.
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
