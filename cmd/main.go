package main

import (
	"book_api/internal/config"
	"book_api/internal/repository/psql"
	"book_api/internal/service"
	"book_api/internal/transport/rest"
	"book_api/pkg/database"
	"book_api/pkg/hash"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPsqlConnection(database.ConnectionConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher("acbd")

	userRepo := psql.NewUserRepository(db)
	userService := service.NewUserService(userRepo, hasher, []byte("my secret"), cfg.Auth.TokenTTL)

	bookRepo := psql.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)

	bookHandler := rest.NewHandler(bookService, userService)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr:    addr,
		Handler: bookHandler.InitRoutes(),
	}

	log.Info("Server started")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
