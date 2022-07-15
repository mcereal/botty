package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mcereal/go-api-server/handler"
)

// InitializeRouter starts gin engine, then it maps the handlers
// to the gin engine and then returns the engine.
func InitializeRouter(e string) *gin.Engine {

	// Check the environment. If its not in development then
	// the gin engine is set to ReleaseMode which doesn't include
	// some things like debug logging.
	if e != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Start the default gin engine.
	r := gin.Default()

	// Create a new handler based off of the Config struct.
	// The new handler maps the initalized gin Default r.
	handler.NewHandler(&handler.Config{
		R: r,
	})

	// return the router with handlers attached
	return r
}
