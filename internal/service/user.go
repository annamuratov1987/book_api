package service

import (
	"book_api/internal/domain"
	"context"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
}

type UserService struct {
	repo   UserRepository
	hasher PasswordHasher
}

func NewUserService(repository UserRepository, hasher PasswordHasher) *UserService {
	return &UserService{
		repo:   repository,
		hasher: hasher,
	}
}

func (s UserService) SignUp(ctx context.Context, input domain.SignUpInput) (int64, error) {
	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     hash,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s UserService) SignIn(ctx context.Context, input domain.SignInInput) error {
	return nil
}
