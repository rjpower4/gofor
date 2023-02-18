package gofor

import (
	"log"
)

type GoforService struct {
	Repository Repository
	Fetcher    ResourceFetcher
}

func New(repoPath string) (GoforService, error) {
	repo, err := newRepositoryFromFile(repoPath)
	if err != nil {
		return GoforService{}, err
	}

	fetcher := NewConcurrentFetcher()

	return GoforService{
		Repository: repo,
		Fetcher:    fetcher,
	}, nil
}

func (s *GoforService) Fetch(names ...string) {
	toFetch := make([]Resource, 0)

	// Check all first so that we only download one if we can download all
	for _, name := range names {
		r, ok := s.Repository.GetByName(name)
		if !ok {
			log.Printf("unknown resource '%s'\n", name)
			return
		} else {
			toFetch = append(toFetch, r)
		}
	}

	for _, r := range toFetch {
		err := s.Fetcher.Enqueue(r)
		if err != nil {
			log.Printf("failed to enqueue resource: %v\n", err)
			return
		}
	}

	_ = s.Fetcher.Wait()
}
