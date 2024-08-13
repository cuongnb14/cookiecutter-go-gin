package pagination

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"

	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/models"
)

func TestSimplePagination_GetData(t *testing.T) {
	configs.GetDB()
	size := 20
	page := 1

	//queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(`page=%d&size=%d`, page, size)

	stmt := configs.GetDB().Model(&models.Agency{})
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	p := &Pagination{}
	p.With(stmt).Request(req)
	pageData := p.Response(&[]models.Agency{})

	//TranslateItems[models.Agency, ShortAgency](pageData)

	log.Println(pageData.Items)
	log.Println(pageData.Total)

}

type ShortAgency struct {
	Apikey string `json:"apikey" binding:"required"`
}
