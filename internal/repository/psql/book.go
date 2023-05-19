package psql

import (
	"book_api/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	if !book.Validate() {
		return 0, domain.ErrorEmptyRequiredField
	}

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
	if err == sql.ErrNoRows {
		return book, domain.ErrorBookNotFound
	}

	return book, err
}

func (r BookRepository) Update(ctx context.Context, id int64, input domain.UpdateBookInput) error {
	fields := make([]string, 0)
	fieldId := 0
	args := make([]interface{}, 0)

	if input.Title != nil {
		fieldId++
		fields = append(fields, fmt.Sprintf("title=$%d", fieldId))
		args = append(args, input.Title)
	}

	if input.Author != nil {
		fieldId++
		fields = append(fields, fmt.Sprintf("author=$%d", fieldId))
		args = append(args, input.Author)
	}

	if input.PublishDate != nil {
		fieldId++
		fields = append(fields, fmt.Sprintf("publish_date=$%d", fieldId))
		args = append(args, input.PublishDate)
	}

	if input.Rating != nil {
		fieldId++
		fields = append(fields, fmt.Sprintf("rating=$%d", fieldId))
		args = append(args, input.Rating)
	}

	if fieldId == 0 {
		return domain.ErrorEmptyUpdateBookInput
	}

	query := fmt.Sprintf("update books set %s where id=%d", strings.Join(fields, ", "), id)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r BookRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec("delete from books where id=$1", id)
	return err
}
