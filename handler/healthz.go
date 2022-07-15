package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz is a handler for the healthz route.
// It sets basic headers and sends a 200 status code with an "OK message"
func (h *Handler) Healthz(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(http.StatusOK, "OK")
}
