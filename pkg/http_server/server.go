package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/net/http2"
)

// Server is a wrapper around the gin engine.
type Server struct {
	*gin.Engine
	config *Config
}

// New creates a new Server.
func New(config *Config) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(otelgin.Middleware(config.OtelServiceName))
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.Use(corsMiddleware(config.Cors))
	engine.Use(requestIDMiddleware())
	engine.Use(secureHeadersMiddleware())

	engine.GET("/metrics", prometheusHandler())
	engine.GET("/healthz", healthzHandler())

	return &Server{
		Engine: engine,
		config: config,
	}, nil
}

// Run starts the server.
func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.Engine,
	}

	if s.config.SSL {
		return srv.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
	}

	http2.ConfigureServer(srv, &http2.Server{})
	return srv.ListenAndServe()
}

// prometheusHandler returns a gin.HandlerFunc for prometheus metrics.
func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// healthzHandler returns a gin.HandlerFunc for health checks.
func healthzHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

// corsMiddleware returns a gin.HandlerFunc for CORS.
func corsMiddleware(config CorsConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.AllowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// requestIDMiddleware returns a gin.HandlerFunc for request IDs.
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-ID", c.GetHeader("X-Request-ID"))
		c.Next()
	}
}

// secureHeadersMiddleware returns a gin.HandlerFunc for secure headers.
func secureHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		if c.Request.TLS != nil {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

func (s *Server) AddGroup(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.Engine.Group(relativePath, handlers...)
}

func (s *Server) AddRoute(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Engine.Handle(httpMethod, relativePath, handlers...)
}

// Config is the server configuration.
type Config struct {
	Port            int        `mapstructure:"port"`
	SSL             bool       `mapstructure:"ssl"`
	CertFile        string     `mapstructure:"cert_file"`
	KeyFile         string     `mapstructure:"key_file"`
	Mode            string     `mapstructure:"mode"`
	OtelServiceName string     `mapstructure:"otel_service_name"`
	Cors            CorsConfig `mapstructure:"cors"`
}

// CorsConfig is the CORS configuration.
type CorsConfig struct {
	AllowedOrigins string `mapstructure:"allowed_origins"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.Mode == "" {
		c.Mode = "debug"
	}
	if c.OtelServiceName == "" {
		c.OtelServiceName = "my-app"
	}
	if c.Cors.AllowedOrigins == "" {
		c.Cors.AllowedOrigins = "*"
	}
	return nil
}

func (s *Server) RunTLS(certFile, keyFile string) error {
	return s.Engine.RunTLS(fmt.Sprintf(":%d", s.config.Port), certFile, keyFile)
}

func (s *Server) SetTrustedProxies(trustedProxies []string) error {
	return s.Engine.SetTrustedProxies(trustedProxies)
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.Engine,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
