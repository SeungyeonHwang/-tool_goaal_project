package myapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPathHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	indexHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatal("Failed!!", res.Code)
	}
}
