package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/synera-br/lockari-backend-app/pkg/logger"
)

// Handler is a generic handler.
type Handler struct {
	log logger.Logger
}

// New creates a new handler.
func New(log logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// Request is a generic request.
type Request struct {
	Payload string `json:"payload"`
}

// Response is a generic response.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// Handle handles a request.
func (h *Handler) Handle(handler func(payload string) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			h.log.Errorf("failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, Response{Error: "invalid request payload"})
			return
		}

		data, err := handler(req.Payload)
		if err != nil {
			h.log.Errorf("failed to handle request: %v", err)
			c.JSON(http.StatusInternalServerError, Response{Error: "failed to handle request"})
			return
		}

		c.JSON(http.StatusOK, Response{Data: data})
	}
}
