package services

import (
	"sync"
)

type DirectoryToProcess struct {
	Path string
}

type Pool struct {
	numberOfThreads int
	tasksChan       chan DirectoryToProcess
	outputChan      chan FileDef
	wg              sync.WaitGroup
}

func (p *Pool) AddTask(t DirectoryToProcess) {
	p.wg.Add(1)
	p.tasksChan <- t
}

func (p *Pool) Run() {
	for i := 0; i < p.numberOfThreads; i++ {
		go func() {
			for task := range p.tasksChan {
				GoOverFiles(task.Path, p)
				p.wg.Done()
			}
		}()
	}
}

func (p *Pool) Close() {
	p.wg.Wait()
	close(p.tasksChan)
}

func NewPool(numberOfThreads int, outputChan chan FileDef) *Pool {
	return &Pool{
		numberOfThreads: numberOfThreads,
		tasksChan:       make(chan DirectoryToProcess, numberOfThreads),
		outputChan:      outputChan,
	}
}
