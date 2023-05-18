package psql

import (
	"book_api/internal/domain"
	"context"
	"database/sql"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r BookRepository) Create(ctx context.Context, book domain.Book) error {
	_, err := r.db.Exec(
		"insert into books (title, author, publish_date, rating) values ($1, $2, $3, $4)",
		book.Title,
		book.Author,
		book.PublishDate,
		book.Rating,
	)
	if err != nil {
		return err
	}

	return nil
}
