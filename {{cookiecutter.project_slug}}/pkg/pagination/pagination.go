package pagination

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func NewStandardPagination() *Pagination {
	return &Pagination{
		maxSize: 50,
	}
}

type Pagination struct {
	maxSize       int
	Query         *gorm.DB
	RequestParams *RequestParams
	withCount     bool
}

type Page struct {
	Items interface{}
	Total int64
}

type RequestParams struct {
	Size int
	Page int
}

func (p *Pagination) With(query *gorm.DB) *Pagination {
	p.Query = query
	return p
}

func (p *Pagination) SetCount(count bool) *Pagination {
	p.withCount = count
	return p
}

func (p *Pagination) Request(request *http.Request) *Pagination {
	query := request.URL.Query()
	requestParams := &RequestParams{}
	if i, e := strconv.Atoi(query.Get("size")); nil == e {
		requestParams.Size = i
	} else {
		requestParams.Size = 50
	}

	if i, e := strconv.Atoi(query.Get("page")); nil == e {
		requestParams.Page = i
	} else {
		requestParams.Page = 1
	}
	p.RequestParams = requestParams
	return p
}

func (p *Pagination) Response(results interface{}) *Page {
	var total int64
	dbs := p.Query.Statement.DB.Session(&gorm.Session{NewDB: true})
	query := dbs.Unscoped().Table("(?) AS s", p.Query)
	if p.withCount {
		query = query.Count(&total).Limit(p.RequestParams.Size).Offset((p.RequestParams.Page - 1) * p.RequestParams.Size)
	} else {
		query = query.Limit(p.RequestParams.Size).Offset((p.RequestParams.Page - 1) * p.RequestParams.Size)
	}
	query.Find(results)

	return &Page{
		Items: results,
		Total: total,
	}
}
