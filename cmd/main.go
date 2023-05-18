package main

import (
	"book_api/internal/repository/psql"
	"book_api/internal/service"
	"book_api/internal/transport/rest"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=123 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	bookRepo := psql.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := rest.NewBookHandler(bookService)

	server := http.Server{
		Addr:    "localhost:8001",
		Handler: bookHandler.InitRoutes(),
	}

	fmt.Println("Server started...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
