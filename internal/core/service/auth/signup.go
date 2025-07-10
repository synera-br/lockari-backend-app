package service

import (
	"context"
	"fmt"
	"log"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	core "github.com/synera-br/lockari-backend-app/internal/core/entity/types"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	"github.com/synera-br/lockari-backend-app/pkg/database"
	"github.com/synera-br/lockari-backend-app/pkg/tokengen"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type SignupEvent struct {
	repo     entity.SignupEventRepository
	auth     authenticator.Authenticator
	tokenJWT tokengen.TokenGenerator
}

func InitializeSignupEventService(repo entity.SignupEventRepository, auth authenticator.Authenticator, tokenJWT tokengen.TokenGenerator) (entity.SignupEventService, error) {

	if repo == nil {
		return nil, core.ErrRepositoryNotFound("SignupEventRepository")
	}

	if auth == nil {
		return nil, core.ErrRepositoryNotFound("Authenticator")
	}

	if tokenJWT == nil {
		return nil, core.ErrRepositoryNotFound("TokenGenerator")
	}

	return &SignupEvent{
		repo:     repo,
		auth:     auth,
		tokenJWT: tokenJWT,
	}, nil
}

func (s *SignupEvent) Create(ctx context.Context, signupData *entity.Signup) (entity.SignupEvent, error) {
	// CHECK STRUCTURE
	fmt.Println("[SERVICE] 1 Creating signup event with data:", signupData)
	if signupData == nil {
		return nil, core.ErrGenericError("Signup event is required")
	}

	if err := signupData.IsValid(); err != nil {
		return nil, fmt.Errorf("Invalid signup event: %w", err)
	}

	// CHECK CONTEXT
	fmt.Println("[SERVICE] 2 Checking context for cancellation")
	if ctx.Err() != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	// CHECK TOKEN
	fmt.Println("[SERVICE] 3 Checking token")
	token := utils.GetTokenFromContext(ctx) // Ensure user ID is retrieved from context

	_, err := s.tokenJWT.Validate(token)
	if err != nil {
		return nil, fmt.Errorf(utils.GenericError, err.Error())
	}

	// CHECK IF USER ALREADY HAS A TENANT
	log.Println("[SERVICE] 4 Checking if user already has a tenant")
	existingTenant, err := s.auth.GetTenant(ctx, token)
	if err == nil && existingTenant != "" {
		return nil, core.ErrGenericError("User already has a tenant assigned")
	}

	// CHECK USER FROM INTERFACE
	fmt.Println("[SERVICE] 5 Checking user from interface")
	tenantId := utils.GenerateTenant()

	// CREATE SIGNUP EVENT
	fmt.Println("[SERVICE] 6 Creating signup event")
	signup := entity.NewSignup(signupData.User, signupData.ClientInfo, tenantId)
	if signup == nil {
		return nil, core.ErrGenericError("Failed to create signup event")
	}

	// VALIDATE SIGNUP EVENT
	fmt.Println("[SERVICE] 7 Validating signup event")
	if err := signup.IsValid(); err != nil {
		return nil, core.ErrGenericError("Invalid signup event")
	}

	// SET TENANT ID
	fmt.Println("[SERVICE] 8 Setting tenant ID")
	err = signup.SetTenant(&tenantId)
	if err != nil {
		return nil, err
	}

	// Log da geração do tenant para auditoria
	log.Printf("Generated tenant %s ", tenantId)

	// SET TENANT ID AT FIREBASE AUTHENTICATION
	fmt.Println("[SERVICE] 9 Setting tenant ID in authentication service")
	err = s.auth.SetTenantId(ctx, signupData.User.Uid, tenantId)
	if err != nil {
		return nil, core.ErrGenericError("Failed to set tenant ID")
	}

	// CONVERT SIGNUP TO MAP
	fmt.Println("[SERVICE] 10 Converting signup data to map")
	data, err := utils.StructToMap(signup.GetSignup())
	if err != nil {
		return nil, core.ErrGenericError("Failed to convert signup data to map")
	}

	// CREATE SIGNUP EVENT
	fmt.Println("[SERVICE] 11 Creating signup event in repository")
	result, err := s.repo.Create(ctx, data)
	if err != nil {
		// Salvar o erro original antes de tentar rollback
		originalErr := err

		// Tentar rollback do tenant
		if rollbackErr := s.auth.SetTenantRollback(ctx, signupData.User.Uid, ""); rollbackErr != nil {
			// Log o erro de rollback, mas retorna o erro original
			log.Printf("Failed to rollback tenant for user %s: %v", signupData.User.Uid, rollbackErr)
		}

		// Sempre retorna o erro original, não o erro do rollback
		return nil, originalErr
	}

	// Log de sucesso para auditoria
	log.Printf("Successfully created signup event for user %s with tenant %s", signupData.User.Uid, tenantId)

	return result, nil
}

func (s *SignupEvent) Get(ctx context.Context, id string) (entity.SignupEvent, error) {

	if ctx.Err() != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	if id == "" {
		return nil, core.ErrGenericError("Signup event ID is required")
	}

	token := utils.GetTokenFromContext(ctx)

	userFromToken, err := s.auth.GetUserID(ctx, token)
	if err != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	// Buscar o signup event primeiro
	filter := database.Conditional{
		Field:  "id",
		Value:  id,
		Filter: database.FilterEquals,
	}

	result, err := s.repo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Verificar se o usuário tem permissão para acessar este signup event
	if result.GetUser().Uid != userFromToken {
		return nil, core.ErrUnauthorized("User is not authorized to access this signup event")
	}

	SignupEvent := entity.NewSignup(result.GetUser(), result.GetClientInfo(), result.GetTenant())
	if err := SignupEvent.IsValid(); err != nil {
		return nil, err
	}

	return SignupEvent, nil
}

func (s *SignupEvent) List(ctx context.Context) ([]entity.SignupEvent, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	token := utils.GetTokenFromContext(ctx)

	userFromToken, err := s.auth.GetUserID(ctx, token)
	if err != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	filter := database.Conditional{
		Field:  "userId",
		Value:  userFromToken,
		Filter: database.FilterEquals,
	}

	result, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, core.ErrNotFound("No signup events found for user")
	}

	var signupEvents []entity.SignupEvent
	for _, signup := range result {
		if err := signup.IsValid(); err != nil {
			return nil, err
		}

		e := entity.NewSignup(signup.GetUser(), signup.GetClientInfo(), signup.GetTenant())
		signupEvents = append(signupEvents, e)
	}

	return signupEvents, nil
}
