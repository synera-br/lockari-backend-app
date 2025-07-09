package repository

import (
	"context"
	"errors"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	"github.com/synera-br/lockari-backend-app/pkg/database"
	"github.com/synera-br/lockari-backend-app/pkg/utils"
)

type LoginEvent struct {
	db database.FirebaseDBInterface
}

func InitializeLoginEventRepository(db database.FirebaseDBInterface) (entity.LoginEventRepository, error) {

	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	if !db.IsConnected() {
		return nil, errors.New("database connection is not initialized")
	}

	return &LoginEvent{
		db: db,
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

func (s *LoginEvent) Get(ctx context.Context, filters []map[string]interface{}) (login entity.LoginEvent, err error) {
	_, err = utils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	if len(filters) == 0 {
		return nil, err
	}

	return login, err
}

func (s *LoginEvent) List(ctx context.Context) (logins []entity.LoginEvent, err error) {
	return logins, err
}
