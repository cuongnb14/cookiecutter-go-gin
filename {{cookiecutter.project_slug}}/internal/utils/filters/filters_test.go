package gormfilter

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestFilterBuilder(t *testing.T) {
	query := fmt.Sprintf(`name=%s&age=1`, "page")

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	b := FilterBuilder{req: req}
	b.AddFilter(Filter{
		Param: "name",
		Field: "name",
		Type:  "string",
		Op:    "=",
	}).AddFilter(Filter{
		Param: "age",
		Field: "age",
		Type:  "int",
		Op:    "=",
	})

	//fmt.Println(b.GetFilterQuery())
}
