package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlessioPani/go-greenlight/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Define an envelope type.
type envelope map[string]any

// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded JSON.
	// Here we use no line prefix ("") and tab indents ("\t") for each element.
	// MashalIndent is quite slower compared to Marshal. For resource-constrained application
	// use json.Marshal instead.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Loop through the header map and add each header to the http.ResponseWriter header map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Define a readJSON() helper for reading requests. This takes the destination
// http.ResponseWriter, the source *http.Request and the target variable dst and check for
// several common errors regarding decoding a JSON (request body).
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding that means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Check for SyntaxError errors.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Check for a generic ErrUnexpectedEOF error.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Check for UnmarshalTypeError errors. These occur when the JSON value is the wrong type
		// for the target destination. If the error relates to a specific field, then we include
		// that in our error message to make it easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Check for EOF error, that occurs when the body request is empty.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Check if the JSON contains a field which cannot be mapped to the target destination.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Check for Max Byte Error that occurs when the request body exceeded our size limit of 1MB.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		// Check for InvalidUnmarshalError errors. These occur when we pass a nil pointer to the
		// Decode function.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Call Decode() again, if the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Retrieve parameters from URL.
	params := httprouter.ParamsFromContext(r.Context())

	// Get ID.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// readString returns a string value from the query string, or the default value if no match
// has been provided.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// Gets a string value from the query.
	s := qs.Get(key)

	// If no key exists, returns the default value.
	if s == "" {
		return defaultValue
	}

	return s
}

// readCSVr reads a string value from the query string and then splits it into a slice on
// the comma character. If no matching key could be found, it returns the provided default value.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Gets a string value from the query.
	csv := qs.Get(key)

	// If no key exists, returns the default value.
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Gets a string value from the query.
	s := qs.Get(key)

	// If no key exists, returns the default value.
	if s == "" {
		return defaultValue
	}

	// Checks for errors during the str -> int conversation. If so, use the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i

}
