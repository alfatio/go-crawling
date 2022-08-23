package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexing(t *testing.T) {

	rc := httptest.NewRecorder()

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/indexing", nil)

	r.ServeHTTP(rc, req)

	if rc.Code != 200 {
		t.Fail()
	}
}
