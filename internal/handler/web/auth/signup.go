package webhandler

import (
	"log"

	"github.com/gin-gonic/gin"
	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	mid "github.com/synera-br/lockari-backend-app/internal/handler/middleware"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	cryptserver "github.com/synera-br/lockari-backend-app/pkg/crypt/crypt_server"
	"github.com/synera-br/lockari-backend-app/pkg/tokengen"
)

type signupHandler struct {
	svc        entity.SignupEventService
	encryptor  cryptserver.CryptDataInterface
	authClient authenticator.Authenticator
	token      tokengen.TokenGenerator
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
	token tokengen.TokenGenerator,
	routerGroup *gin.RouterGroup,
	middleware ...gin.HandlerFunc,
) SignupHandlerInterface {
	handler := &signupHandler{
		svc:        svc,
		encryptor:  encryptData,
		authClient: authClient,
		token:      token,
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

	withoutAuth := routerGroup.Group("/")
	withoutAuth.POST("/api/v1/audit/auth", h.WithJWT)
	withoutAuth.POST("/audit/auth", h.WithJWT)
	withoutAuth.GET("/api/v1/audit/auth", h.WithJWT)
	withoutAuth.GET("/audit/auth", h.WithJWT)

	// extra := routerGroup.Group("/")
	// for _, mw := range middleware {
	// 	extra.Use(mw)
	// }
	// extra.GET("/api/v1/audit/auth", h.Extras)
	// extra.GET("/audit/auth", h.Extras)

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

	var body cryptserver.CryptData
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	log.Println("Received payload:", body)
	decryptedData, err := h.encryptor.PayloadData(body.Payload)
	if err != nil {
		log.Println("Error decrypting payload:", err)
		c.JSON(400, gin.H{"error": "Error processing request data"})
		return
	}

	log.Println("Decrypted data:", string(decryptedData))

	c.JSON(200, gin.H{"message": "Extras handled successfully"})
}

func (h *signupHandler) WithJWT(c *gin.Context) {
	log.Println("Handling signup with JWT...")

	// Exemplo de geração de token normal
	claims := tokengen.TokenClaims{
		AppID: "lockari-frontend-app",
		Scope: []string{"signup", "read"},
		Metadata: map[string]interface{}{
			"type": "signup",
		},
	}

	nonExpiringToken, err := h.token.GenerateNonExpiring(claims)
	if err != nil {
		log.Printf("Error generating non-expiring token: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate non-expiring token"})
		return
	}

	validatedClaims, err := h.token.Validate(nonExpiringToken)
	if err != nil {
		log.Printf("Non-expiring token validation failed: %v", err)
	} else {
		log.Printf("Non-expiring token valid for user: %s, NonExpiring: %v",
			validatedClaims.UserID, validatedClaims.NonExpiring)
	}

	// Verificar se tokens expiram
	log.Printf("Non-expiring token expired: %v", h.token.IsExpired(nonExpiringToken))
	log.Println("Non-expiring token:", nonExpiringToken)

	c.JSON(200, gin.H{
		"message": "JWT tokens generated successfully",
		"tokens": map[string]interface{}{
			"non_expiring": nonExpiringToken,
		},
	})
}
