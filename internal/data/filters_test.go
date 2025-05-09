package data

import (
	"testing"

	"github.com/AlessioPani/go-greenlight/internal/validator"
)

var sortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

// Test method used to test the private function calculateMetadata.
func TestCalculateMetadata(t *testing.T) {
	// Tests to be run.
	tests := []struct {
		name           string
		totalRecords   int
		page           int
		pageSize       int
		expectedResult Metadata
	}{
		{"all positive int", 10, 2, 2, Metadata{CurrentPage: 2, PageSize: 2, FirstPage: 1, LastPage: 5, TotalRecords: 10}},
		{"totalRecords = 0", 0, 2, 2, Metadata{}},
		{"page = 0", 10, 0, 2, Metadata{}},
		{"pageSize = 0", 10, 2, 0, Metadata{}},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateMetadata(test.totalRecords, test.page, test.pageSize)
			if result != test.expectedResult {
				t.Errorf("got %+v, expected %+v", result, test.expectedResult)
			}
		})
	}
}

// Test method used to test the ValidateFilters function.
func TestValidateFilters(t *testing.T) {
	// Tests to be run.
	tests := []struct {
		name           string
		filter         Filters
		field          string
		expectedResult string
	}{
		{"filter ok", Filters{Page: 1, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, "", ""},
		{"page not greater than 0", Filters{Page: 0, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, "page", "must be greater than zero"},
		{"page greater than 10.000.000", Filters{Page: 10_000_001, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, "page", "must be a maximum of 10 million"},
		{"page size not greater than 0", Filters{Page: 1, PageSize: 0, Sort: "title", SortSafeList: sortSafeList}, "page_size", "must be greater than zero"},
		{"page size greater than 100", Filters{Page: 1, PageSize: 101, Sort: "title", SortSafeList: sortSafeList}, "page_size", "must be a maximum of 100"},
		{"filter not supported", Filters{Page: 1, PageSize: 2, Sort: "wrong", SortSafeList: sortSafeList}, "sort", "invalid sort value"},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := validator.New()
			ValidateFilters(v, test.filter)
			if v.Errors[test.field] != test.expectedResult {
				t.Errorf("got %s, expected %s", v.Errors[test.field], test.expectedResult)
			}
		})
	}
}

// Test method used to test the sortColumn function.
func TestSortColumn(t *testing.T) {
	// Tests to be run.
	tests := []struct {
		name           string
		filter         Filters
		expectedResult string
	}{
		{"no hyphen", Filters{Page: 1, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, "title"},
		{"hyphen", Filters{Page: 1, PageSize: 2, Sort: "-title", SortSafeList: sortSafeList}, "title"},
		{"panic", Filters{Page: 1, PageSize: 2, Sort: "wrong", SortSafeList: sortSafeList}, ""},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Defer function to verify if the sortColumn function did panic during
			// the "panic" test execution.
			defer func() {
				if err := recover(); test.name == "panic" && err == nil {
					t.Errorf("the function did not panic")
				}
			}()

			// Checks the results.
			result := test.filter.sortColumn()
			if result != test.expectedResult {
				t.Errorf("got %s, expected %s", test.filter.sortColumn(), test.expectedResult)
			}
		})
	}
}

// Test method used to test the sortDirection function.
func TestSortDirection(t *testing.T) {
	// Tests to be run.
	tests := []struct {
		name           string
		filter         Filters
		expectedResult string
	}{
		{"asc", Filters{Page: 1, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, "ASC"},
		{"desc", Filters{Page: 1, PageSize: 2, Sort: "-title", SortSafeList: sortSafeList}, "DESC"},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.filter.sortDirection()
			if result != test.expectedResult {
				t.Errorf("got %s, expected %s", result, test.expectedResult)
			}
		})
	}
}

// Test method used to test the limit function.
func TestLimit(t *testing.T) {
	// Tests to be run.
	tests := []struct {
		name           string
		filter         Filters
		expectedResult int
	}{
		{"pagesize 2", Filters{Page: 1, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, 2},
		{"pagesize 3", Filters{Page: 1, PageSize: 3, Sort: "-title", SortSafeList: sortSafeList}, 3},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.filter.limit()
			if result != test.expectedResult {
				t.Errorf("got %d, expected %d", result, test.expectedResult)
			}
		})
	}
}

// Test method used to test the offset function.
func TestOffset(t *testing.T) {
	// Tests to be run.
	tests := []struct {
		name           string
		filter         Filters
		expectedResult int
	}{
		{"page 1", Filters{Page: 1, PageSize: 2, Sort: "title", SortSafeList: sortSafeList}, 0},
		{"page 2", Filters{Page: 2, PageSize: 3, Sort: "-title", SortSafeList: sortSafeList}, 3},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.filter.offset()
			if result != test.expectedResult {
				t.Errorf("got %d, expected %d", result, test.expectedResult)
			}
		})
	}
}
