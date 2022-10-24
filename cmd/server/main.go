package main

import (
	"fmt"

	"github.com/marcos-nsantos/aluraflix-api/internal/database"
)

func run() error {
	fmt.Println("Starting up the application")

	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err = db.Migrate("file:///migrations"); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
