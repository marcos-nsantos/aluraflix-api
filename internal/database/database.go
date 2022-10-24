package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	Client *sqlx.DB
}

func Connect() (*Database, error) {
	var count uint
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	for {
		db, err := sqlx.Connect("postgres", connectionString)
		if err != nil {
			log.Println("Postgres is not ready yet, retrying...")
			count++
		} else {
			return &Database{Client: db}, nil
		}

		if count > 10 {
			return nil, err
		}

		log.Println("Waiting for 5 seconds to retry...")
		time.Sleep(5 * time.Second)
		continue
	}
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
