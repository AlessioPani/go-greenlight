package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Declare a custom Runtime type, which has the underlying type int32 (the same as our
// Movie struct field).
type Runtime int32

// Define an error that our UnmarshalJSON() method can return if we're unable to parse
// or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}

// Implement a UnmarshalJSON() method on the Runtime type so that it satisfies the
// json.Unmarshaler interface.
// IMPORTANT: Because UnmarshalJSON() needs to modify the receiver (our Runtime type),
// we must use a pointer receiver for this to work correctly.
// Otherwise, we will only be modifying a copy (which is then discarded when this method returns).
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// We expect to receive a field like "<runtime> mins", so first of all
	// we remove the quotes "".
	// If no quotes are available, we raise an invalid runtime format error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the input to separate the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanitize input by checking its format.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Parse the number.
	num, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// If the input is in the valid format, assign
	// the number to the receiver.
	*r = Runtime(num)

	return nil
}
