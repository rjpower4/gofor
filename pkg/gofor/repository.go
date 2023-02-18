package gofor

import (
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
)

// Repository is the general interface for types that support retrieving resources
type Repository interface {
	// GetByName returns the given resource from the repository
	GetByName(name string) (Resource, bool)

	// GetAll returns all resources from the repository
	GetAll() []Resource

	// Insert inserts a resource into the repository
	Insert(name string, resource Resource) bool

	// Remove removes a resource from the repository
	Remove(name string) bool

	// Update changes the resource with the given name to the given resource
	Update(name string, resource Resource) bool
}

type repository struct {
	Resources map[string]Resource `toml:"resources"`
}

func (r repository) GetByName(name string) (Resource, bool) {
	val, ok := r.Resources[name]
	return val, ok
}

func (r repository) GetAll() []Resource {
	resources := make([]Resource, 0)
	for _, value := range r.Resources {
		resources = append(resources, value)
	}
	return resources
}

func (r repository) Insert(name string, resource Resource) bool {
	// Don't do anything if it exists in the repository
	_, ok := r.Resources[name]
	if ok {
		return false
	}

	r.Resources[name] = resource
	return true
}

func (r repository) Remove(name string) bool {
	_, ok := r.Resources[name]

	if ok {
		delete(r.Resources, name)
	}

	return ok
}

func (r repository) Update(name string, resource Resource) bool {
	_, ok := r.Resources[name]

	if ok {
		r.Resources[name] = resource
	}

	return ok
}

func newRepositoryFromFile(filepath string) (Repository, error) {

	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("unable to read file: %v\n", err)
		return &repository{}, err
	}

	var r repository
	err = toml.Unmarshal(fileBytes, &r)
	if err != nil {
		log.Printf("unable to unmarshal file: %v\n", err)
		return &repository{}, err
	}

	return &r, nil
}
