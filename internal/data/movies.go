package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/validator"
	"github.com/lib/pq"
)

// Struct that defines a movie object.
type Movie struct {
	ID        int64     `json:"id"`                // Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`                 // Timestamp for when the movie is added to our database
	Title     string    `json:"title"`             // Movie title
	Year      int32     `json:"year,omitempty"`    // Movie release year
	Runtime   Runtime   `json:"runtime,omitempty"` // Movie runtime (in minutes)
	Genres    []string  `json:"genres,omitempty"`  // Slices of genres of the movie (drama, comedy, etc..)
	Version   int32     `json:"version"`           // The version number starts at 1 and will be incremented each time the movie information is updated
}

// Movie model struct that wraps a db connection pool.
type MovieModel struct {
	DB *sql.DB
}

// Insert is a method for inserting a new record in the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	// SQL query for inserting a movie in the db and returning
	// the system-generated data.
	query := `INSERT INTO movies (title, year, runtime, genres)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, version`

	// Values for the placeholders in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Executes QueryRow in order to get the system-generated data and returns the error, if any.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get is a method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	// Sanitize ID.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// SQL query for getting a movie in the db by its ID.
	query := `SELECT id, created_at, title, year, runtime, genres, version
			  FROM movies
			  WHERE id = $1`

	// movie struct to be returned back.
	movie := Movie{}

	// Executes QueryRow in order to get the record.
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// Checks for errors during Scan().
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

// Update is a method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// Delete is a method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return nil
}

// ValidateMovie is a method that validates input data.
func ValidateMovie(v *validator.Validator, movie *Movie) {
	// Title
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	// Year
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// Runtime
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	// Genres
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must be provided at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must contain only unique values")
}
