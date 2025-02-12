package data

import (
	"database/sql"
	"errors"
)

var (
	// ErrRecordNotFound is a custom error returned if the MovieModel.Get() method
	// has as input an Id of a movie that doesn't exist in our database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is a custom error returned if a data race condition
	// occurred while editing a record.
	ErrEditConflict = errors.New("edit confict")
)

// Models is a struct which wraps our models.
type Models struct {
	Movies      MovieModelInterface
	Permissions PermissionModelInterface
	Tokens      TokenModelInterface
	Users       UserModelInterface
}

// NewModels() method returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      &MovieModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      &TokenModel{DB: db},
		Users:       &UserModel{DB: db},
	}
}
