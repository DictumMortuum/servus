package boardgames

import (
	"github.com/DictumMortuum/servus/pkg/generic"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"id":1,"name":"boardgameprices"}`

	r, err := generic.SetupMainRouter("../../servusrc")
	if err != nil {
		t.Error(err)
	}

	rest := r.Group("/rest/v1")
	rest.GET("/store/:id", generic.F(GetStore))
	rest.GET("/store", generic.F(GetListStore))

	req, _ := http.NewRequest("GET", "/rest/v1/store/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	if mockResponse != string(responseData) {
		t.Errorf("%s\n", string(responseData))
	}
}
