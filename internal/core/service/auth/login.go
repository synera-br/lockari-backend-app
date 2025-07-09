package service

import (
	"context"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	core "github.com/synera-br/lockari-backend-app/internal/core/entity/types"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type LoginEvent struct {
	repo entity.LoginEventRepository
}

func InitializeLoginEventService(repo entity.LoginEventRepository) (entity.LoginEventService, error) {

	if repo == nil {
		return nil, core.ErrRepositoryNotFound("LoginEventRepository")
	}

	return &LoginEvent{
		repo: repo,
	}, nil
}

func (s *LoginEvent) Create(ctx context.Context, requestLogin entity.LoginEvent) (login entity.LoginEvent, err error) {

	uid, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	if *uid != requestLogin.GetUser().Uid {
		return nil, err
	}

	if err = requestLogin.IsValid(); err != nil {
		return nil, err
	}

	return requestLogin, err
}

func (s *LoginEvent) Get(ctx context.Context, id string) (login entity.LoginEvent, err error) {
	_, err = utils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	return login, err
}

func (s *LoginEvent) List(ctx context.Context) (logins []entity.LoginEvent, err error) {
	return logins, err
}
