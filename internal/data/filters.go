package data

import (
	"strings"

	"github.com/AlessioPani/go-greenlight/internal/validator"
)

// Filters is a struct that contains filters parameter of a request query.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// sortColumn checks if a column matches one of the entries of the safe list, and returns
// the string without the hyphen character, if exists.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter:" + f.Sort)
}

// sortDirection returns the sort direction depending on the prefix character of the Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// limit is a method that returns the pagesize of a filter.
func (f Filters) limit() int {
	return f.PageSize
}

// offset is a method that returns the page of a filter.
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// ValidateFilters is an helper method to validate a Filters struct.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// Metadata is a struct that contains pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calculateMetadata is a function that calculates the pagination metadata,
// given the total number of records, current page and the page size.
func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
