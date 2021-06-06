package goapi

import (
	"net/url"
	"strings"
)

// Sorting ...
func (config *Config) Sorting(urlQuery *url.Values) (sortOrder string, sortBy string) {
	validSortOrders := map[string]bool{
		"desc": true,
		"asc":  true,
	}
	sortOrder = strings.ToLower(urlQuery.Get("sort_order"))

	// Set default sort order and sort by
	if !validSortOrders[sortOrder] {
		sortOrder = config.Sort.DefaultSortOrder
	}

	sortBy = strings.ToLower(urlQuery.Get("sort_by"))
	if !SliceContains(config.Sort.ValidSortBy, sortBy) {
		sortBy = config.Sort.DefaultSortBy
	}

	return
}
