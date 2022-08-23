package db

import (
	"database/sql"
	"log"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func DB() *sql.DB {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DB_CONNECTION_STRING"))

	if err != nil {
		panic(err)
	}

	return db
}
