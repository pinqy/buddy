package dbmanager

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var client *DBClient

// var TABLE_NAMES = []string{"category", "tag", "expense"}

// TODO: get rid of hard coded creds
var cfg = mysql.Config{
	User:   "root",
	Passwd: "Blues5/20",
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "buddy",
}

func initDB() (*sql.DB, error) {
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	// Ping test db connection
	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	log.Printf("Successfully connected to buddy database!")
	return db, nil
}

func GetDBClient() (*DBClient, error) {
	if client != nil {
		return client, nil
	}

	db, err := initDB()
	if err != nil {
		return nil, err
	}
	return &DBClient{database: db}, nil
}
