package rest

import (
	"book_api/internal/domain"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type BookService interface {
	Create(ctx context.Context, book domain.Book) error
}

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(books BookService) BookHandler {
	return BookHandler{
		bookService: books,
	}
}

func (h *BookHandler) InitRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/create", h.createBook).Methods(http.MethodPost)

	return r
}

func (h BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("BookHandler.createBook() read body error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book domain.Book
	err = json.Unmarshal(reqBody, &book)
	if err != nil {
		log.Printf("BookHandler.createBook() json unmarshal error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.bookService.Create(r.Context(), book)
	if err != nil {
		log.Printf("BookHandler.createBook() create book error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
