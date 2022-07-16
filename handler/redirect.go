package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Redirect is a handler to redirect to GitHub Repo.
func (h *Handler) Redirect(c *gin.Context) {
	c.Redirect(http.StatusFound, "https://github.com/mcereal/botty")
}
