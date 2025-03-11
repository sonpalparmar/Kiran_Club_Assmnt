// internal/repository/job_repository.go
package repository

import (
	"errors"
	"sync"

	"kirana-club-assignment/backend/internal/models"
)

var (
	ErrJobNotFound = errors.New("job not found")
)

type JobRepository interface {
	Create(job models.Job) (string, error)
	GetByID(id string) (models.Job, error)
	Update(job models.Job) error
}

type InMemoryJobRepository struct {
	mu   sync.RWMutex
	jobs map[string]models.Job
}

func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{
		jobs: make(map[string]models.Job),
	}
}

func (r *InMemoryJobRepository) Create(job models.Job) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return job.ID, nil
}

func (r *InMemoryJobRepository) GetByID(id string) (models.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, exists := r.jobs[id]
	if !exists {
		return models.Job{}, ErrJobNotFound
	}
	return job, nil
}

func (r *InMemoryJobRepository) Update(job models.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.jobs[job.ID]; !exists {
		return ErrJobNotFound
	}
	r.jobs[job.ID] = job
	return nil
}
