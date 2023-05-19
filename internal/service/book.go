package service

import (
	"book_api/internal/domain"
	"context"
)

type BookRepository interface {
	Create(ctx context.Context, book domain.Book) (int64, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetById(ctx context.Context, id int64) (domain.Book, error)
	Update(ctx context.Context, id int64, input domain.UpdateBookInput) error
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s BookService) Create(ctx context.Context, book domain.Book) (int64, error) {
	return s.repo.Create(ctx, book)
}

func (s BookService) GetAll(ctx context.Context) ([]domain.Book, error) {
	return s.repo.GetAll(ctx)
}

func (s BookService) GetById(ctx context.Context, id int64) (domain.Book, error) {
	return s.repo.GetById(ctx, id)
}

func (s BookService) Update(ctx context.Context, id int64, input domain.UpdateBookInput) error {
	return s.repo.Update(ctx, id, input)
}
