package webhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/synera-br/lockari-backend-app/internal/core/service/auth"
	"github.com/synera-br/lockari-backend-app/internal/handler/web"
)

// SignupHandler is the signup handler.
type SignupHandler struct {
	authService *auth.Service
	handler     *web.Handler
}

// NewSignupHandler creates a new signup handler.
func NewSignupHandler(
	authService *auth.Service,
	handler *web.Handler,
) *SignupHandler {
	return &SignupHandler{
		authService: authService,
		handler:     handler,
	}
}

// RegisterRoutes registers the signup routes.
func (h *SignupHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/signup", h.handler.Handle(h.signup))
}

func (h *SignupHandler) signup(payload string) (interface{}, error) {
	return nil, h.authService.Signup(nil, payload)
}
