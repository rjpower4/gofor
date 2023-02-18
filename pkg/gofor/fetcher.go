package gofor

import (
	"context"
	"golang.org/x/sync/semaphore"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

type ResourceFetcher interface {
	// Enqueue queues up the resource for fetching
	Enqueue(r Resource) error

	// Cancel all currently running fetches
	Cancel()

	// Wait blocks until all enqueued resource fetches have completed
	Wait() error
}

type ConcurrentFetcher struct {
	WorkerCount int64
	client      http.Client
	ctx         context.Context
	semaphore   *semaphore.Weighted
}

func NewConcurrentFetcher() *ConcurrentFetcher {
	nWorkers := int64(runtime.NumCPU())
	return &ConcurrentFetcher{
		WorkerCount: nWorkers,
		client:      http.Client{},
		ctx:         context.TODO(),
		semaphore:   semaphore.NewWeighted(nWorkers),
	}
}

func (cf *ConcurrentFetcher) Enqueue(r Resource) error {
	err := cf.semaphore.Acquire(cf.ctx, 1)
	if err != nil {
		log.Printf("unable to acquire semaphore: %v\n", err)
		return err
	}

	go cf.fetch(r.URL, r.FileName)

	return nil
}

func (cf *ConcurrentFetcher) Cancel() {
	panic("Not implemented")
}

func (cf *ConcurrentFetcher) Wait() error {
	err := cf.semaphore.Acquire(cf.ctx, cf.WorkerCount)
	if err != nil {
		log.Printf("failed to acquire semaphore: %v\n", err)
		return err
	}
	return nil
}

func (cf *ConcurrentFetcher) fetch(url string, outPath string) {
	defer cf.semaphore.Release(1)

	f, err := os.Create(outPath)
	if err != nil {
		log.Printf("unable to create file at %s\n", outPath)
		return
	}
	defer f.Close()

	response, err := cf.client.Get(url)
	if err != nil {
		log.Printf("unable to perform HTTP request: %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("bad status: %d\n", response.StatusCode)
		return
	}

	_, err = io.Copy(f, response.Body)
	if err != nil {
		log.Printf("unable to copy data to file: %v\n", err)
	}
}
