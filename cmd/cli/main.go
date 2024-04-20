package main

import (
	"database/sql"
	"duplicates/internal/repositories"
	"duplicates/internal/services"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	folder := "."
	if len(os.Args) > 1 {
		folder = os.Args[1]
	}

	db, err := sql.Open("sqlite3", "./duplicates.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	repositories.Create(db)

	outputChan := make(chan services.FileDef)
	go func() {
		for file := range outputChan {
			repositories.InsertFile(db, file)
		}
	}()

	p := services.NewPool(10, outputChan)
	p.Run()
	p.AddTask(services.DirectoryToProcess{Path: folder})
	p.Close()
}
