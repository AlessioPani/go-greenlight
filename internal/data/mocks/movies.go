package mocks

import (
	"errors"
	"strings"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/data"
)

var ValidMovie = data.Movie{
	ID:        1,
	CreatedAt: time.Now(),
	Title:     "Lord of the Rings",
	Year:      2001,
	Runtime:   228,
	Genres:    []string{"fantasy", "epic"},
	Version:   1,
}

var InvalidMovie = data.Movie{
	ID:        2,
	CreatedAt: time.Now(),
	Title:     "",
	Year:      2001,
	Runtime:   228,
	Genres:    []string{"fantasy", "epic"},
	Version:   1,
}

var ErrorMovie = data.Movie{
	ID:        3,
	CreatedAt: time.Now(),
	Title:     "Watchmen",
	Year:      2009,
	Runtime:   162,
	Genres:    []string{"action"},
	Version:   1,
}

var metadata = data.Metadata{
	CurrentPage:  1,
	PageSize:     2,
	FirstPage:    1,
	LastPage:     1,
	TotalRecords: 1,
}

type MovieModel struct{}

// Method used for mocking the Insert method for the Movie model.
func (m *MovieModel) Insert(movie *data.Movie) error {
	if movie.Title == ErrorMovie.Title {
		return errors.New("db error")
	}

	return nil
}

// Method used for mocking the Get method for the Movie model.
func (m *MovieModel) Get(id int64) (*data.Movie, error) {
	switch id {
	case 1:
		return &ValidMovie, nil
	case 2:
		return &InvalidMovie, nil
	case 3:
		return &ErrorMovie, nil
	case 4:
		return nil, data.ErrRecordNotFound
	default:
		return nil, errors.New("server error")
	}
}

// Method used for mocking the GetAll method for the Movie model.
func (m *MovieModel) GetAll(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
	if strings.Contains(strings.ToLower(ErrorMovie.Title), strings.ToLower(title)) {
		return nil, data.Metadata{}, errors.New("server error")
	}

	return []*data.Movie{&ValidMovie}, metadata, nil
}

// Method used for mocking the Update method for the Movie model.
func (m *MovieModel) Update(movie *data.Movie) error {
	if movie.ID == ErrorMovie.ID {
		return data.ErrEditConflict
	}

	return nil
}

// Method used for mocking the Delete method for the Movie model.
func (m *MovieModel) Delete(id int64) error {
	if id == 4 {
		return data.ErrRecordNotFound
	}

	return nil
}
