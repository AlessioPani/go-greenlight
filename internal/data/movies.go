package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/validator"
	"github.com/lib/pq"
)

// Movie represents a movie record.
type Movie struct {
	ID        int64     `json:"id"`                // Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`                 // Timestamp for when the movie is added to our database
	Title     string    `json:"title"`             // Movie title
	Year      int32     `json:"year,omitempty"`    // Movie release year
	Runtime   Runtime   `json:"runtime,omitempty"` // Movie runtime (in minutes)
	Genres    []string  `json:"genres,omitempty"`  // Slices of genres of the movie (drama, comedy, etc..)
	Version   int32     `json:"version"`           // The version number starts at 1 and will be incremented each time the movie information is updated
}

// MovieModelInterface defines the movie data operations.
type MovieModelInterface interface {
	Insert(ctx context.Context, movie *Movie) error
	Get(ctx context.Context, id int64) (*Movie, error)
	GetAll(ctx context.Context, title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	Update(ctx context.Context, movie *Movie) error
	Delete(ctx context.Context, id int64) error
}

// MovieModel provides access to movie records in the database.
type MovieModel struct {
	DB *sql.DB
}

// Insert is a method for inserting a new record in the movies table.
func (m *MovieModel) Insert(parentContext context.Context, movie *Movie) error {
	// SQL query for inserting a movie in the db and returning
	// the system-generated data.
	query := `INSERT INTO movies (title, year, runtime, genres)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, version`

	// Values for the placeholders in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Create a context with a three-second timeout.
	ctx, cancel := context.WithTimeout(parentContext, 3*time.Second)
	defer cancel()

	// Executes QueryRow in order to get the system-generated data and returns the error, if any.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get is a method for fetching a specific record from the movies table.
func (m *MovieModel) Get(parentContext context.Context, id int64) (*Movie, error) {
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

	// Create a context with a three-second timeout.
	ctx, cancel := context.WithTimeout(parentContext, 3*time.Second)
	defer cancel()

	// Executes QueryRow in order to get the record.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// Handle errors returned while scanning the row.
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

// GetAll retrieves movies matching the supplied filters.
func (m *MovieModel) GetAll(parentContext context.Context, title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	// SQL query to retrieve the requested records.
	// Uses built-in functionality of Postgres to achieve full-text search with lexemes.
	sortColumn, err := filters.sortColumn()
	if err != nil {
		return nil, Metadata{}, err
	}
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
			 	          FROM movies
			  			  WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
			  			  AND (genres @> $2 OR $2 ='{}')
			  			  ORDER BY %s %s, id ASC
						  LIMIT $3 OFFSET $4`, sortColumn, filters.sortDirection())

	// Create a context with a three-second timeout.
	ctx, cancel := context.WithTimeout(parentContext, 3*time.Second)
	defer cancel()

	// Executes the query.
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres), filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	// Scan rows to fill the result and metadata struct.
	totalRecords := 0
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	// Check for errors encountered while iterating over the rows.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}

// Update is a method for updating a specific record in the movies table.
//
// This method has an optimistic locking functionality to prevent data race
// condition on an update for a specific record.
// It updates a record only if the version of the record is the same at
// the beginning of the request, otherwise returns an ErrEditConflict error.
func (m *MovieModel) Update(parentContext context.Context, movie *Movie) error {
	// SQL query for updating a movie record.
	query := `UPDATE movies
			  SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			  WHERE id = $5 AND version = $6
			  RETURNING version`

	// Values for the placeholders in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version}

	// Create a context with a three-second timeout.
	ctx, cancel := context.WithTimeout(parentContext, 3*time.Second)
	defer cancel()

	// Executes QueryRow in order to get the system-generated data and returns the error, if any.
	// If no row has been retrieved, return an Edit Conflict (data race condition) error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete is a method for deleting a specific record from the movies table.
func (m *MovieModel) Delete(parentContext context.Context, id int64) error {
	// Sanitize ID.
	if id < 1 {
		return ErrRecordNotFound
	}

	// SQL query for deleting a movie.
	query := `DELETE FROM movies
			  WHERE id = $1`

	// Create a context with a three-second timeout.
	ctx, cancel := context.WithTimeout(parentContext, 3*time.Second)
	defer cancel()

	// Executes the query.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// ValidateMovie checks that the movie fields contain valid values.
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
