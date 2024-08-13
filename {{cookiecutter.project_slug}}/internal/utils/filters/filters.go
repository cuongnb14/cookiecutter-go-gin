package gormfilter

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"{{ cookiecutter.project_slug }}/configs"
)

const (
	FTypeString  = "string"
	FTypeUuid    = "uuid"
	FTypeDecimal = "decimal"
	FTypeInt     = "int"
)

type Filter struct {
	Param string
	Field string
	Type  string
	Op    string
}

func (f *Filter) GetQuery(value string) (string, interface{}, error) {
	var where string
	var params interface{}
	switch f.Type {
	case FTypeString:
		params = value
	case FTypeUuid:
		v, err := uuid.FromBytes([]byte(value))
		if err != nil {
			return "", nil, err
		}
		params = v
	case FTypeInt:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", nil, err
		}
		params = v
	case FTypeDecimal:
		v, err := decimal.NewFromString(value)
		if err != nil {
			return "", nil, err
		}
		params = v
	}

	where = fmt.Sprintf("%s %s ?", f.Field, f.Op)
	return where, params, nil
}

type FilterBuilder struct {
	Filters  []Filter
	req      *http.Request
	Ordering []string

	Query *gorm.DB
}

func (b *FilterBuilder) SetQuery(query *gorm.DB) *FilterBuilder {
	b.Query = query
	return b
}

func (b *FilterBuilder) SetOrdering(ordering []string) *FilterBuilder {
	b.Ordering = ordering
	return b
}

func (b *FilterBuilder) SetRequest(req *http.Request) *FilterBuilder {
	b.req = req
	return b
}

func (b *FilterBuilder) AddFilter(filter Filter) *FilterBuilder {
	b.Filters = append(b.Filters, filter)
	return b
}

func (b *FilterBuilder) BuildFilter() (string, []interface{}, string) {
	logger := configs.GetLogger()
	var wheres []string
	var params []interface{}
	query := b.req.URL.Query()
	for _, filter := range b.Filters {
		valueParam := query.Get(filter.Param)
		if valueParam == "" {
			continue
		}
		w, p, err := filter.GetQuery(valueParam)
		if err != nil {
			logger.Error(err)
		} else {
			wheres = append(wheres, w)
			params = append(params, p)
		}
	}

	ordering := ""

	if len(b.Ordering) > 0 {
		orderingParam := query.Get("order")
		orderingType := query.Get("order_type")
		if orderingParam != "" {
			if slices.Contains(b.Ordering, orderingParam) {
				ordering = fmt.Sprintf("%s %s", orderingParam, orderingType)
			}
		}
	}

	return strings.Join(wheres, " and "), params, ordering
}

func (b *FilterBuilder) GetFilterQuery() *gorm.DB {
	wheres, params, ordering := b.BuildFilter()
	if wheres != "" {
		b.Query = b.Query.Where(wheres, params...)
	}
	if ordering != "" {
		b.Query = b.Query.Order(ordering)
	}
	return b.Query
}
