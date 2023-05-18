package main

import (
	"book_api/internal/repository/psql"
	"book_api/internal/service"
	"book_api/internal/transport/rest"
	"book_api/pkg/database"
	"log"
	"net/http"
)

func main() {
	connConfig := database.ConnectionConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Username: "postgres",
		Password: "123",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	db, err := database.NewPsqlConnection(connConfig)
	if err != nil {
		log.Fatal(err)
	}

	bookRepo := psql.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := rest.NewBookHandler(bookService)

	server := http.Server{
		Addr:    "localhost:8001",
		Handler: bookHandler.InitRoutes(),
	}

	log.Println("Server started...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
