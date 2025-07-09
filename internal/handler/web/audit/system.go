package webhandler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/audit"
	"github.com/synera-br/lockari-backend-app/internal/handler/middleware"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	cryptserver "github.com/synera-br/lockari-backend-app/pkg/crypt/crypt_server"
	"github.com/synera-br/lockari-backend-app/pkg/tokengen"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type auditSystemEventHandler struct {
	svc        entity.AuditSystemEventService
	encryptor  cryptserver.CryptDataInterface
	authClient authenticator.Authenticator
	token      tokengen.TokenGenerator
}

type auditSystemEventHandlerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

func InicializeAuditSystemEventHandler(svc entity.AuditSystemEventService, encryptor cryptserver.CryptDataInterface, authClient authenticator.Authenticator, token tokengen.TokenGenerator, routerGroup *gin.RouterGroup, middlewares ...gin.HandlerFunc) (auditSystemEventHandlerInterface, error) {

	if svc == nil {
		return nil, fmt.Errorf(utils.ServiceNotFoundError, "audit system event service")
	}

	if encryptor == nil {
		return nil, fmt.Errorf(utils.ServiceNotFoundError, "audit system event encryptor")
	}

	if authClient == nil {
		return nil, fmt.Errorf(utils.ServiceNotFoundError, "audit system event auth client")
	}

	if token == nil {
		return nil, fmt.Errorf(utils.ServiceNotFoundError, "audit system event token generator")
	}

	handler := &auditSystemEventHandler{
		svc:        svc,
		encryptor:  encryptor,
		authClient: authClient,
		token:      token,
	}

	handler.setupRoutes(routerGroup, middlewares...)

	return handler, nil
}

func (h *auditSystemEventHandler) setupRoutes(routerGroup *gin.RouterGroup, middlewares ...gin.HandlerFunc) {

	auditRoutes := routerGroup.Group("/audit")
	middlewares = append(middlewares, middleware.ValidateToken(&gin.Context{}, h.authClient))
	for _, mw := range middlewares {
		auditRoutes.Use(mw)
	}

	auditRoutes.POST("/auth", h.Create)
	auditRoutes.GET("/auth", h.List)
	auditRoutes.GET("/auth/:id", h.Get)

}

func (h *auditSystemEventHandler) Create(c *gin.Context) {

	token := c.GetHeader("X-Token")
	dataToken, err := h.authClient.ValidateToken(c.Request.Context(), token)
	if err != nil {
		log.Println("Error validating token:", err)
		c.JSON(401, gin.H{"error": "Invalid or expired token"})
		return
	}

	log.Println("Received token data:", dataToken)
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
	c.JSON(200, gin.H{"message": "Audit event created successfully"})
}

func (w *auditSystemEventHandler) Get(c *gin.Context) {

}

func (w *auditSystemEventHandler) List(c *gin.Context) {

}
