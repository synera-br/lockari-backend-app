package service

import (
	"context"
	"fmt"
	"log"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	core "github.com/synera-br/lockari-backend-app/internal/core/entity/types"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	"github.com/synera-br/lockari-backend-app/pkg/database"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type SignupEvent struct {
	repo entity.SignupEventRepository
	auth authenticator.Authenticator
}

func InitializeSignupEventService(repo entity.SignupEventRepository, auth authenticator.Authenticator) (entity.SignupEventService, error) {

	if repo == nil {
		return nil, core.ErrRepositoryNotFound("SignupEventRepository")
	}

	if auth == nil {
		return nil, core.ErrRepositoryNotFound("Authenticator")
	}

	return &SignupEvent{
		repo: repo,
		auth: auth,
	}, nil
}

func (s *SignupEvent) Create(ctx context.Context, signup entity.SignupEvent) (entity.SignupEvent, error) {
	if err := signup.IsValid(); err != nil {
		return nil, core.ErrGenericError("Invalid signup event")
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	token, err := utils.GetTokenFromContext(ctx, s.auth) // Ensure user ID is retrieved from context
	if err != nil {
		return nil, err
	}

	user := signup.GetUser()
	userFromToken, err := s.auth.GetUserID(ctx, token)
	if err != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	if userFromToken != user.Uid {
		return nil, core.ErrUnauthorized("User ID does not match signup event user ID")
	}

	// Verificar se o usuário já possui um tenant
	existingTenant, err := s.auth.GetTenant(ctx, token)
	if err == nil && existingTenant != "" {
		return nil, core.ErrGenericError("User already has a tenant assigned")
	}

	tenantId := utils.GenerateTenant()

	// Log da geração do tenant para auditoria
	log.Printf("Generated tenant %s for user %s", tenantId, userFromToken)

	err = signup.SetTenant(&tenantId)
	if err != nil {
		return nil, err
	}

	err = s.auth.SetTenantId(ctx, userFromToken, tenantId)
	if err != nil {
		return nil, core.ErrGenericError("Failed to set tenant ID")
	}

	data, err := utils.StructToMap(signup.GetSignup())
	if err != nil {
		return nil, core.ErrGenericError("Failed to convert signup data to map")
	}

	// Tentar criar no repository
	result, err := s.repo.Create(ctx, data)
	if err != nil {
		// Salvar o erro original antes de tentar rollback
		originalErr := err

		// Tentar rollback do tenant
		if rollbackErr := s.auth.SetTenantRollback(ctx, userFromToken, ""); rollbackErr != nil {
			// Log o erro de rollback, mas retorna o erro original
			log.Printf("Failed to rollback tenant for user %s: %v", userFromToken, rollbackErr)
		}

		// Sempre retorna o erro original, não o erro do rollback
		return nil, originalErr
	}

	// Log de sucesso para auditoria
	log.Printf("Successfully created signup event for user %s with tenant %s", userFromToken, tenantId)

	return result, nil
}

func (s *SignupEvent) Get(ctx context.Context, id string) (entity.SignupEvent, error) {

	if ctx.Err() != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	if id == "" {
		return nil, core.ErrGenericError("Signup event ID is required")
	}

	token, err := utils.GetTokenFromContext(ctx, s.auth)
	if err != nil {
		return nil, err
	}

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

	return result, nil
}

func (s *SignupEvent) List(ctx context.Context) ([]entity.SignupEvent, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	token, err := utils.GetTokenFromContext(ctx, s.auth)
	if err != nil {
		return nil, err
	}

	userFromToken, err := s.auth.GetUserID(ctx, token)
	if err != nil {
		return nil, fmt.Errorf(utils.ContextCancelled, ctx.Err().Error())
	}

	filter := database.Conditional{
		Field:  "userId",
		Value:  userFromToken,
		Filter: database.FilterEquals,
	}

	return s.repo.List(ctx, filter)
}
