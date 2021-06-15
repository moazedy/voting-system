package httpEngine

import (
	"math"
	"strconv"
	"voting-system/constants"
	"voting-system/helper"

	"github.com/gin-gonic/gin"
)

type httpPagination struct {
	helper.Pagination
}

func newHttpPagination(pagination helper.Pagination) httpPagination {
	return httpPagination{Pagination: pagination}
}

// From gin context get page, per_page, last_id and order from gin context
func (p *httpPagination) FromGinContext(c *gin.Context) *helper.Paging {
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", string(constants.PaginationDefaultPerPage)))
	if err != nil {
		perPage = constants.PaginationDefaultPerPage
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	lastId := c.DefaultQuery("last_id", "")
	order := c.DefaultQuery("order", constants.PaginationDefaultOrder)
	// construct pagination type
	paging := helper.Paging{
		PerPage: int(math.Min(float64(perPage), float64(p.GetMaxPerPage()))),
		Page:    page,
		Order:   order,
		LastID:  lastId,
	}
	// Set maximum records for a page
	paging.SetMaxPerPage(p.GetMaxPerPage())
	return &paging
}
