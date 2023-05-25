package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

var ErrorUserNotFound = errors.New("user with such credentials not found")

func init() {
	validate = validator.New()
}

type User struct {
	ID           int64
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type SignUpInput struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (inp SignUpInput) Validate() error {
	return validate.Struct(inp)
}

func (inp SignInInput) Validate() error {
	return validate.Struct(inp)
}
