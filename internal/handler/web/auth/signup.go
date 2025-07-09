package webhandler

import (
	"log"

	"github.com/gin-gonic/gin"
	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	mid "github.com/synera-br/lockari-backend-app/internal/handler/middleware"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	cryptserver "github.com/synera-br/lockari-backend-app/pkg/crypt/crypt_server"
)

type signupHandler struct {
	svc        entity.SignupEventService
	encryptor  cryptserver.CryptDataInterface
	authClient authenticator.Authenticator
}

type SignupHandlerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Extras(c *gin.Context)
}

func InitializeSignupHandler(
	svc entity.SignupEventService,
	encryptData cryptserver.CryptDataInterface,
	authClient authenticator.Authenticator,
	routerGroup *gin.RouterGroup,
	middleware ...gin.HandlerFunc,
) SignupHandlerInterface {
	handler := &signupHandler{
		svc:        svc,
		encryptor:  encryptData,
		authClient: authClient,
	}

	handler.setupRoutes(routerGroup, middleware...)
	return handler
}

func (h *signupHandler) setupRoutes(routerGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) {

	signupRoutes := routerGroup.Group("/auth/signup")
	middleware = append(middleware, mid.ValidateToken(&gin.Context{}, h.authClient))
	for _, mw := range middleware {
		signupRoutes.Use(mw)
	}

	signupRoutes.POST("", h.Create)
	signupRoutes.GET("", h.List)
	signupRoutes.GET("/:id", h.Get)

	extra := routerGroup.Group("/")
	for _, mw := range middleware {
		extra.Use(mw)
	}
	extra.GET("/api/v1/audit/auth", h.Extras)
	extra.GET("/audit/auth", h.Extras)
	extra.POST("/api/v1/audit/auth", h.Extras)
	extra.POST("/audit/auth", h.Extras)

}

func (h *signupHandler) Create(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Signup event created successfully"})
}

func (h *signupHandler) Get(c *gin.Context) {
	log.Println("Retrieving signup event...")
	c.JSON(200, gin.H{"message": "Signup event retrieved successfully"})
}

func (h *signupHandler) List(c *gin.Context) {
	log.Println("Listing signup events...")
	c.JSON(200, gin.H{"message": "Signup events listed successfully"})
}

func (h *signupHandler) Extras(c *gin.Context) {
	log.Println("Handling extras for signup...")
	log.Println(c.Request.Header)
	c.JSON(200, gin.H{"message": "Extras handled successfully"})
}
