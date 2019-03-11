package main_test

import (
	"net/http"
	"testing"
)

func TestGetOriginalURL(t *testing.T) {
	// make dummy req
	response, err := http.Get("http://localhost:8000/v1/short/1")
	if http.StatusOK != response.StatusCode {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.StatusCode)
	}
	if err != nil {
		t.Error("Encountered an error:", err)
	}
}
