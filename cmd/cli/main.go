package main

import (
	"database/sql"
	"duplicates/internal/models"
	"duplicates/internal/repositories"
	"duplicates/internal/utils"
	"fmt"
	"os"
	"path/filepath"

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

	dirPool := utils.NewPool(10)
	outputChan := make(chan models.FileDef)
	go func() {
		for file := range outputChan {
			if file.IsDir {
				dirPool.AddTask(utils.DirectoryToProcess{Path: filepath.Join(file.Path, file.Filename), ResultChan: outputChan})
			} else {
				repositories.InsertFile(db, file)
			}
		}
	}()

	dirPool.Run()
	dirPool.AddTask(utils.DirectoryToProcess{Path: folder, ResultChan: outputChan})
	dirPool.Close()

	fmt.Println("Calculating hashes")
	filePool := utils.NewPool(10)
	outputChan = make(chan models.FileDef)
	go func() {
		for file := range outputChan {
			repositories.UpdateHash(db, file.ID, file.Hash)
		}
	}()

	filePool.Run()
	files := repositories.GetPotentialDuplicates(db)
	for _, f := range files {
		filePool.AddTask(utils.FileToProcess{File: f, ResultChan: outputChan})
	}

	fmt.Println("Found duplicates:")
	duplicates := repositories.GetDuplicates(db)
	for _, f := range duplicates {
		fmt.Println(f.Path, f.Filename)
	}
}
