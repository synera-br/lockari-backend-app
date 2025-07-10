package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU.",
})

var routerPath *gin.RouterGroup

func init() {
	prometheus.MustRegister(cpuTemp)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *RestAPI) Run(handle http.Handler) error {

	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", s.Config.Port),
		Handler: s.Routes.Handler(),
	}

	s.Routes.Use(s.CorsMiddleware())
	http2.ConfigureServer(&srv, &http2.Server{})

	return srv.ListenAndServe()
}

func (s *RestAPI) RunTLS() error {
	return nil
}

func setHeader(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	c.Next()
}

func (s *RestAPI) MiddlewareHeader(c *gin.Context) {
	log.Println("MiddlewareHeader called")
	log.Println("Request Headers:", c.Request.Header)
	authorization := c.GetHeader("X-AUTHORIZATION")
	token := c.GetHeader("X-TOKEN")
	app := c.GetHeader("X-APP")

	if authorization == "" && token == "" {
		c.Errors = append(c.Errors, &gin.Error{
			Err:  fmt.Errorf("missing header parameters for authentication"),
			Type: gin.ErrorTypePublic,
		})
		c.AbortWithStatus(401)
		return
	}

	if app == "" {
		c.Errors = append(c.Errors, &gin.Error{
			Err:  fmt.Errorf("missing header parameters for app"),
			Type: gin.ErrorTypePublic,
		})
		c.AbortWithStatus(401)
		return
	}
	c.Next()
}

func (s *RestAPI) CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", " X-Authorization", " X-USERID", " X-APP", " X-USER", " X-TRACE-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
