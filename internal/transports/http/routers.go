package http

import (
	"go-boilerplate/internal/configs"
	"go-boilerplate/internal/services"
	"go-boilerplate/internal/transports/http/handlers"
	middewares "go-boilerplate/internal/transports/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all API routes grouped by version/module
func RegisterRoutes(r *gin.Engine, svcs services.Register, cfg configs.Config) {
	// inisiate ExampleHandler with the ExampleService from services.Register
	// This allows the handler to use the service for business logic operations.
	// The handler methods will call the service methods to perform actions like creating, updating, or deleting examples.
	// This approach promotes separation of concerns and makes the code more maintainable.
	exampleHandler := handlers.NewExampleHandler(svcs.ExampleService)
	// add more handlers if needed

	// Build middleware from config
	basicAuthMiddleware := middewares.BasicAuthMiddleware(cfg)
	// add more middewares if needed

	// Define the routes for the example module
	exampleRoute := r.Group("/example", basicAuthMiddleware)
	{
		// Users
		exampleRoute.POST("/", exampleHandler.CreateExample)
	}
}
