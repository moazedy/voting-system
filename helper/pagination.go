package helper

import (
	"strings"
)

// Pagination interface, any adapter must implement this interface
type Pagination interface {
	// Setters
	SetMaxPerPage(int)
	// Getters
	GetMaxPerPage() int
	GetNextPage() int
	GetPrevPage() int
	GetOffset() int
	GetLimit() int
	GetOrder() string
	GetLastID() string
	GetPage() int
	// HasNextPage is a lazy loading helper function that help to calculate if results have next page to be call or not
	HasNextPage(int) bool
}

// Paging structure
type Paging struct {
	maxPerPage  int
	PerPage     int
	Page        int
	Order       string
	LastID      string
	hasNextPage bool
}

// NewPagination is pagination constructor
func NewPagination(maxPerPage int) Pagination {
	return &Paging{
		maxPerPage: maxPerPage,
		Page:       1,
	}
}

// SetMaxPerPage set maximum record number of a page
func (p *Paging) SetMaxPerPage(i int) {
	p.maxPerPage = i
}

// GetMaxPerPage is a maxPerPage getter function
func (p *Paging) GetMaxPerPage() int {
	return p.maxPerPage
}

func (p *Paging) GetNextPage() int {
	if !p.hasNextPage {
		return 0
	}
	return p.Page + 1
}

func (p *Paging) GetPrevPage() int {
	if p.Page <= 1 {
		return 0
	}
	return p.Page - 1
}

// GetOffset calculate offset from page and per page values
func (p *Paging) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns per page value
func (p *Paging) GetLimit() int {
	return p.PerPage
}

//GetOrder return order in UPPERCASE
func (p *Paging) GetOrder() string {
	return strings.ToUpper(p.Order)
}

// GetLastID returns last set id
func (p *Paging) GetLastID() string {
	return p.LastID
}

func (p *Paging) HasNextPage(length int) bool {
	if length > p.PerPage {
		p.hasNextPage = true
		return true
	}
	return false
}

func (p *Paging) GetPage() int {
	return p.Page
}
