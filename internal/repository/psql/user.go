package psql

import (
	"book_api/internal/domain"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	result := r.db.QueryRow(
		"insert into users (name, email, password, registered_at) values ($1, $2, $3, $4) returning id",
		user.Name,
		user.Email,
		user.Password,
		user.RegisteredAt,
	)

	var id int64
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r UserRepository) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	row := r.db.QueryRow(
		"select id, name, email, registered_at from users where email=$1 and password=$2", email, password)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)
	if err == sql.ErrNoRows {
		return user, domain.ErrorUserNotFound
	}

	return user, err
}
