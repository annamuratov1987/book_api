package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type ConnectionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPsqlConnection(cnf ConnectionConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cnf.Host, cnf.Port, cnf.Username, cnf.Password, cnf.DBName, cnf.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("database.NewPsqlConnection() connection open error: %s", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("database.NewPsqlConnection() ping error: %s", err)
		return nil, err
	}

	return db, nil
}
