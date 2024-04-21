package utils

import (
	"duplicates/internal/models"
	"duplicates/internal/services"
	"sync"
)

type Task interface {
	Run(p *Pool)
}

// @TODO introduce a callback function to handle the processor and make it dynamic to isolate dependencies
type DirectoryToProcess struct {
	Path       string
	ResultChan chan models.FileDef
}

func (d DirectoryToProcess) Run(p *Pool) {
	services.GoOverFiles(d.Path, d.ResultChan)
}

type FileToProcess struct {
	File       models.FileDef
	ResultChan chan models.FileDef
}

func (f FileToProcess) Run(p *Pool) {
	services.CalculateHash(f.File, f.ResultChan)
}

type Pool struct {
	numberOfThreads int
	tasksChan       chan Task
	wg              sync.WaitGroup
}

func (p *Pool) AddTask(t Task) {
	p.wg.Add(1)
	p.tasksChan <- t
}

func (p *Pool) Run() {
	for i := 0; i < p.numberOfThreads; i++ {
		go func() {
			for task := range p.tasksChan {
				task.Run(p)
				p.wg.Done()
			}
		}()
	}
}

func (p *Pool) Close() {
	p.wg.Wait()
	close(p.tasksChan)
}

func NewPool(numberOfThreads int) *Pool {
	return &Pool{
		numberOfThreads: numberOfThreads,
		tasksChan:       make(chan Task, numberOfThreads),
	}
}
