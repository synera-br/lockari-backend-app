package service

import (
	"context"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	core "github.com/synera-br/lockari-backend-app/internal/core/entity/types"
	"github.com/synera-br/lockari-backend-app/pkg/database"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type SignupEvent struct {
	repo entity.SignupEventRepository
}

func InitializeSignupEventService(repo entity.SignupEventRepository) (entity.SignupEventService, error) {

	if repo == nil {
		return nil, core.ErrRepositoryNotFound("SignupEventRepository")
	}

	return &SignupEvent{
		repo: repo,
	}, nil
}

func (s *SignupEvent) Create(ctx context.Context, signup entity.SignupEvent) (entity.SignupEvent, error) {
	if err := signup.IsValid(); err != nil {
		return nil, core.ErrGenericError("Invalid signup event")
	}

	userId, err := utils.GetUserID(ctx) // Ensure user ID is retrieved from context
	if err != nil {
		return nil, err
	}

	user := signup.GetUser()
	if *userId != user.Uid {
		return nil, core.ErrUnauthorized("User ID does not match signup event user ID")
	}

	tenantId, err := utils.GenerateTenantID()
	if err != nil {
		return nil, core.ErrGenericError("Failed to generate tenant ID")
	}

	err = signup.SetTenant(&tenantId)
	if err != nil {
		return nil, err
	}

	data, err := utils.StructToMap(signup.GetSignup())
	if err != nil {
		return nil, core.ErrGenericError("Failed to convert signup data to map")
	}

	return s.repo.Create(ctx, data)
}

func (s *SignupEvent) Get(ctx context.Context, id string) (entity.SignupEvent, error) {
	userId, err := utils.GetUserID(ctx) // Ensure user ID is retrieved from context
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, core.ErrGenericError("Signup event ID is required")
	}

	if *userId != id {
		return nil, core.ErrUnauthorized("User ID does not match signup event ID")
	}

	filter := database.Conditional{
		Field:  "id",
		Value:  id,
		Filter: database.FilterEquals,
	}

	return s.repo.Get(ctx, filter)
}

func (s *SignupEvent) List(ctx context.Context) ([]entity.SignupEvent, error) {
	userID, err := utils.GetUserID(ctx) // Ensure user ID is retrieved from context
	if err != nil {
		return nil, err
	}

	filter := database.Conditional{
		Field:  "userId",
		Value:  userID,
		Filter: database.FilterEquals,
	}

	return s.repo.List(ctx, filter)
}
