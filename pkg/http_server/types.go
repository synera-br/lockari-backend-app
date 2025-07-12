package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

// SwaggerInfo is the swagger information.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "My App",
	Description:      "This is a sample server.",
	InfoInstanceName: "swagger",
}

// AddSwagger registers the swagger routes.
func (s *Server) AddSwagger(path string) {
	s.GET(fmt.Sprintf("/%s/*any", path), ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// AddUIRoutes registers the UI routes.
func (s *Server) AddUIRoutes(path string) {
	s.GET(fmt.Sprintf("/%s", path), func(c *gin.Context) {
		c.Redirect(301, fmt.Sprintf("/%s/index.html", path))
	})
	s.Static(fmt.Sprintf("/%s", path), "./ui/dist")
}
