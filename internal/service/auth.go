package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/tinhnguyen-git/health-memory-go/internal/model"
	"github.com/tinhnguyen-git/health-memory-go/internal/repo"
)

type AuthService struct {
	users *repo.UserRepo
}

func NewAuthService(u *repo.UserRepo) *AuthService {
	return &AuthService{users: u}
}

func (s *AuthService) Register(email, password, name string) (*model.User, error) {
	_, err := s.users.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email exists")
	}
	// if other error not found continue
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        email,
		PasswordHash: string(pw),
		Name:         name,
		Provider:     "local",
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Authenticate(email, password string) (*model.User, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Provider != "local" {
		return nil, errors.New("account is social provider")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

// Social upsert
func (s *AuthService) UpsertSocial(ctx context.Context, provider, providerID, email, name, avatar string) (*model.User, error) {
	// try find by providerID
	user, err := s.users.FindByProvider(provider, providerID)
	if err == nil {
		return user, nil
	}
	// try find by email to link
	existing, err2 := s.users.FindByEmail(email)
	if err2 == nil {
		existing.Provider = provider
		existing.ProviderID = providerID
		if name != "" {
			existing.Name = name
		}
		if avatar != "" {
			existing.AvatarURL = avatar
		}
		_ = s.users.Update(existing)
		return existing, nil
	}
	// create new
	u := &model.User{
		Email:      email,
		Provider:   provider,
		ProviderID: providerID,
		Name:       name,
		AvatarURL:  avatar,
	}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}
