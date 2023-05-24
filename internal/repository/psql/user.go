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
