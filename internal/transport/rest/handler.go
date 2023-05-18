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
	Create(ctx context.Context, book domain.Book) (int64, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
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
	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllBooks).Methods(http.MethodGet)
	}

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

	lastInsertedId, err := h.bookService.Create(r.Context(), book)
	if err != nil {
		log.Printf("BookHandler.createBook() create book error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(lastInsertedId)
	if err != nil {
		log.Printf("BookHandler.getAllBooks() lastInsertId json marshal error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h BookHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookService.GetAll(r.Context())
	if err != nil {
		log.Printf("BookHandler.getAllBooks() get []book error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	result, err := json.Marshal(books)
	if err != nil {
		log.Printf("BookHandler.getAllBooks() []book json marshal error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}
