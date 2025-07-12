package auth

import (
	"context"
	"encoding/json"

	entity "github.com/synera-br/lockari-backend-app/internal/core/entity/auth"
	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	cryptserver "github.com/synera-br/lockari-backend-app/pkg/crypt/crypt_server"
	"github.com/synera-br/lockari-backend-app/pkg/logger"
	"github.com/synera-br/lockari-backend-app/pkg/tokengen"
)

// Service is the auth service.
type Service struct {
	loginService  entity.LoginEventService
	signupService entity.SignupEventService
	encryptor     cryptserver.CryptDataInterface
	authClient    authenticator.Authenticator
	token         tokengen.TokenGenerator
	log           logger.Logger
}

// New creates a new auth service.
func New(
	loginService entity.LoginEventService,
	signupService entity.SignupEventService,
	encryptor cryptserver.CryptDataInterface,
	authClient authenticator.Authenticator,
	token tokengen.TokenGenerator,
	log logger.Logger,
) *Service {
	return &Service{
		loginService:  loginService,
		signupService: signupService,
		encryptor:     encryptor,
		authClient:    authClient,
		token:         token,
		log:           log,
	}
}

// Login performs a login.
func (s *Service) Login(ctx context.Context, payload string) error {
	decryptedData, err := s.encryptor.PayloadData(payload)
	if err != nil {
		s.log.Errorf("failed to decrypt payload: %v", err)
		return err
	}

	var loginEvent entity.LoginEvent
	if err := json.Unmarshal(decryptedData, &loginEvent); err != nil {
		s.log.Errorf("failed to unmarshal login event: %v", err)
		return err
	}

	if err := loginEvent.IsValid(); err != nil {
		s.log.Errorf("invalid login event: %v", err)
		return err
	}

	// The login logic is not implemented yet.
	return nil
}

// Signup performs a signup.
func (s *Service) Signup(ctx context.Context, payload string) error {
	decryptedData, err := s.encryptor.PayloadData(payload)
	if err != nil {
		s.log.Errorf("failed to decrypt payload: %v", err)
		return err
	}

	var signup entity.Signup
	if err := json.Unmarshal(decryptedData, &signup); err != nil {
		s.log.Errorf("failed to unmarshal signup: %v", err)
		return err
	}

	if _, err := s.signupService.Create(ctx, &signup); err != nil {
		s.log.Errorf("failed to create signup: %v", err)
		return err
	}

	return nil
}
