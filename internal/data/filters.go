package data

import "github.com/AlessioPani/go-greenlight/internal/validator"

// Filters is a struct that contains filters parameter of a request query.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// ValidateFilters is an helper method to validate a Filters struct.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
