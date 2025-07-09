package service

import (
	"context"
	"errors"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/audit"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type auditSystemEvent struct {
	repo entity.AuditSystemEventRepository
	auth authenticator.Authenticator
}

func InitializeAuditSystemEventService(repo entity.AuditSystemEventRepository, auth authenticator.Authenticator) (entity.AuditSystemEventService, error) {
	if repo == nil {
		return nil, errors.New(utils.RepositoryNotFound + "AuditSystemEventRepository")
	}

	if auth == nil {
		return nil, errors.New(utils.RepositoryNotFound + "Authenticator")
	}

	return &auditSystemEvent{
		repo: repo,
		auth: auth,
	}, nil
}

func (s *auditSystemEvent) Create(ctx context.Context, event *entity.AuditSystemEvent) (*entity.AuditSystemEvent, error) {

	// CONTEXT
	if ctx.Err() != nil {
		return nil, errors.New(utils.ContextCancelled)
	}
	token, err := utils.GetTokenFromContext(ctx, s.auth) // Ensure user ID is retrieved from context
	if err != nil {
		return nil, err
	}

	if token == "" {
		return nil, errors.New(utils.InvalidToken)
	}

	// AUDIT
	if err := event.IsValid(); err != nil {
		return nil, err
	}

	data, err := utils.StructToMap(event)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *auditSystemEvent) Get(ctx context.Context, id string) (*entity.AuditSystemEvent, error) {
	return nil, nil
}

func (a *auditSystemEvent) List(ctx context.Context) ([]entity.AuditSystemEvent, error) {
	return nil, nil
}
