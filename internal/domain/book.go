package domain

import (
	"errors"
	"time"
)

var (
	ErrorEmptyRequiredField   = errors.New("required field is empty")
	ErrorBookNotFound         = errors.New("book not found")
	ErrorEmptyUpdateBookInput = errors.New("empty update book input")
)

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int       `json:"rating"`
}

type UpdateBookInput struct {
	Title       *string    `json:"title"`
	Author      *string    `json:"author"`
	PublishDate *time.Time `json:"publish_date"`
	Rating      *int       `json:"rating"`
}

func (b Book) Validate() bool {
	if b.Title == "" || b.Author == "" {
		return false
	}

	return true
}
