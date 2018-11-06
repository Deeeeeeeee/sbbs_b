package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserPage(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	router.ServeHTTP(w, req)

	println(w.Body.String())
}
