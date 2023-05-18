package service

import (
	"book_api/internal/domain"
	"context"
)

type BookRepository interface {
	Create(ctx context.Context, book domain.Book) error
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s BookService) Create(ctx context.Context, book domain.Book) error {
	err := s.repo.Create(ctx, book)
	if err != nil {
		return err
	}

	return nil
}
