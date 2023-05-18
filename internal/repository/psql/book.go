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

func (r BookRepository) Create(ctx context.Context, book domain.Book) (int64, error) {
	result := r.db.QueryRow(
		"insert into books (title, author, publish_date, rating) values ($1, $2, $3, $4) returning id",
		book.Title,
		book.Author,
		book.PublishDate,
		book.Rating,
	)

	var id int64
	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r BookRepository) GetAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := r.db.Query("select id, title, author, publish_date, rating from books")
	if err != nil {
		return nil, err
	}

	books := make([]domain.Book, 0)
	for rows.Next() {
		var book domain.Book

		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (r BookRepository) GetById(ctx context.Context, id int64) (domain.Book, error) {
	row := r.db.QueryRow("select id, title, author, publish_date, rating from books where id=$1", id)

	var book domain.Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
	if err != nil {
		return domain.Book{}, err
	}

	return book, err
}
