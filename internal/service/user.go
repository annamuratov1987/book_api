package service

import (
	"book_api/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type UserService struct {
	repo       UserRepository
	hasher     PasswordHasher
	hmacSecret []byte
	tokenTTL   time.Duration
}

func NewUserService(repository UserRepository, hasher PasswordHasher, secret []byte, tokenTTL time.Duration) *UserService {
	return &UserService{
		repo:       repository,
		hasher:     hasher,
		hmacSecret: secret,
		tokenTTL:   tokenTTL,
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

func (s UserService) SignIn(ctx context.Context, input domain.SignInInput) (string, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", domain.ErrorUserNotFound
		}
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": strconv.FormatInt(user.ID, 10),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(s.tokenTTL).Unix(),
	})

	return token.SignedString(s.hmacSecret)
}

func (s UserService) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["user"].(string)
	if !ok {
		return 0, errors.New("invalid user")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid user id")
	}

	return int64(id), nil
}
