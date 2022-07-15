package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHealthz checks the /healthz endpoint
func TestHealthz(t *testing.T) {
	w := httptest.NewRecorder()
	r := gin.Default()
	h := &Handler{}
	r.GET("/healthz", h.Healthz)

	r.ServeHTTP(w, httptest.NewRequest("GET", "/healthz", nil))

	if w.Code != http.StatusOK {
		t.Error("Did not get expect HTTP status code, got", w.Code)
	}
	if w.Body.String() != "OK" {
		t.Error("Did not get expected Status message, got", w.Body.String())
	}
}
