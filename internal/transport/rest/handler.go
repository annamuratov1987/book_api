package rest

import (
	"book_api/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type BookService interface {
	Create(ctx context.Context, book domain.Book) (int64, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetById(ctx context.Context, id int64) (domain.Book, error)
	Update(ctx context.Context, id int64, in domain.UpdateBookInput) error
	Delete(ctx context.Context, id int64) error
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
	r.Use(requestLogging)

	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getBookById).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods(http.MethodPut)
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods(http.MethodDelete)
	}

	return r
}

func (h BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "reading request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book domain.Book
	err = json.Unmarshal(reqBody, &book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "unmarshal request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lastInsertedId, err := h.bookService.Create(r.Context(), book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "service error",
		}).Error(err)

		if errors.Is(err, domain.ErrorEmptyRequiredField) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(lastInsertedId)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "lastInsertId json marshal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h BookHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookService.GetAll(r.Context())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getAllBooks",
			"problem": "get []book error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(books)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getAllBooks",
			"problem": "[]book json marshal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}

func (h BookHandler) getBookById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"problem": "get id from request error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetById(r.Context(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"problem": "service error",
		}).Error(err)

		if errors.Is(err, domain.ErrorBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookById",
			"problem": "book json marshal error",
		}).Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}

func (h BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "get id from request error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetById(r.Context(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "service error",
		}).Error(err)

		if errors.Is(err, domain.ErrorBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "read request body error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input domain.UpdateBookInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "json unmarshal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.bookService.Update(r.Context(), book.ID, input)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "service error",
		}).Error(err)

		if errors.Is(err, domain.ErrorEmptyUpdateBookInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "get id from request error",
		}).Error(err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetById(r.Context(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "service error",
		}).Error(err)

		if errors.Is(err, domain.ErrorBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = h.bookService.Delete(r.Context(), book.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "service error",
		}).Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
