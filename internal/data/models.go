package data

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound is a custom error returned if the MovieModel.Get() method
// has as input an Id of a movie that doesn't exist in our database.
var ErrRecordNotFound = errors.New("record not found")

// Models is a struct which wraps our models.
type Models struct {
	Movies MovieModel
}

// NewModels() method returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
