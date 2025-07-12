package webhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/synera-br/lockari-backend-app/internal/core/service/auth"
	"github.com/synera-br/lockari-backend-app/internal/handler/web"
)

// LoginHandler is the login handler.
type LoginHandler struct {
	authService *auth.Service
	handler     *web.Handler
}

// NewLoginHandler creates a new login handler.
func NewLoginHandler(
	authService *auth.Service,
	handler *web.Handler,
) *LoginHandler {
	return &LoginHandler{
		authService: authService,
		handler:     handler,
	}
}

// RegisterRoutes registers the login routes.
func (h *LoginHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/login", h.handler.Handle(h.login))
}

func (h *LoginHandler) login(payload string) (interface{}, error) {
	return nil, h.authService.Login(nil, payload)
}
