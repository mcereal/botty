package handler

import (
	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type Handler struct{}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R *gin.Engine
}

// NewHandler creates a new handler from the Handler struct that takes in
// the gin engine from Config. It maps routes to handler functions.
func NewHandler(c *Config) {
	h := &Handler{} // currently has no properties

	c.R.GET("/", h.Redirect)

	// create a router group to handle path prefix
	g := c.R.Group("api/v1")
	g.POST("/githubwebhook/payload", h.GitHubWebhookHandler)

	// base route without router group prefix
	c.R.GET("/healthz", h.Healthz)
}
