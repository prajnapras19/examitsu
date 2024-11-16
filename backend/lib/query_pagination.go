package lib

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"gorm.io/gorm"
)

type QueryPagination struct {
	Page     int
	PageSize int
	Sort     string // note: fill this field only in backend
}

func (p *QueryPagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *QueryPagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *QueryPagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *QueryPagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id ASC"
	}
	return p.Sort
}

func (p *QueryPagination) Scope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort())
	}
}

func GetQueryPaginationFromContext(c *gin.Context) (*QueryPagination, error) {
	res := QueryPagination{}

	page, err := strconv.ParseUint(
		c.DefaultQuery(constants.QueryParameterPage, constants.DefaultValueQueryParameterPage),
		10,
		64,
	)
	if err != nil {
		return nil, err
	}
	res.Page = int(page)

	pageSize, err := strconv.ParseUint(
		c.DefaultQuery(constants.QueryParameterPageSize, constants.DefaultValueQueryParameterPageSize),
		10,
		64,
	)
	if err != nil {
		return nil, err
	}
	res.PageSize = int(pageSize)

	return &res, nil
}

func GetDefaultPagination() *QueryPagination {
	return &QueryPagination{
		Page:     constants.DefaultQueryPaginationPage,
		PageSize: constants.DefaultQueryPaginationPageSize,
	}
}
