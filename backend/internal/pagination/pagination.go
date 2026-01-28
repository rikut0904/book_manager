package pagination

import (
	"net/http"
	"strconv"
	"strings"

	"book_manager/backend/internal/config"
)

// Params holds pagination parameters.
type Params struct {
	Page     int
	PageSize int
}

// ParseParams extracts pagination parameters from a request.
// Uses defaultPageSize if not specified in the query string.
func ParseParams(r *http.Request, defaultPageSize int) Params {
	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	if page > config.MaxPageNumber {
		page = config.MaxPageNumber
	}

	pageSize := parsePositiveInt(r.URL.Query().Get("pageSize"), defaultPageSize)
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}
}

// SliceRange calculates start and end indices for slicing a list.
// Returns (start, end) where end is exclusive.
func (p Params) SliceRange(totalLen int) (start, end int) {
	start = (p.Page - 1) * p.PageSize
	if start >= totalLen {
		return totalLen, totalLen
	}
	end = start + p.PageSize
	if end > totalLen {
		end = totalLen
	}
	return start, end
}

func parsePositiveInt(value string, defaultVal int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return defaultVal
	}
	return parsed
}
