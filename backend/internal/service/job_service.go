// internal/service/job_service.go
package service

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"kirana-club-assignment/backend/internal/models"
	"kirana-club-assignment/backend/internal/processor"
	"kirana-club-assignment/backend/internal/repository"
)

type JobService interface {
	CreateJob(req models.JobRequest) (string, error)
	GetJobStatus(jobID string) (models.JobStatusResponse, error)
}

type jobService struct {
	jobRepo        repository.JobRepository
	storeRepo      repository.StoreRepository
	imageProcessor *processor.ImageProcessor
}

func NewJobService(
	jobRepo repository.JobRepository,
	storeRepo repository.StoreRepository,
	imageProcessor *processor.ImageProcessor,
) JobService {
	return &jobService{
		jobRepo:        jobRepo,
		storeRepo:      storeRepo,
		imageProcessor: imageProcessor,
	}
}

func (s *jobService) CreateJob(req models.JobRequest) (string, error) {
	// Validate request
	if req.Count != len(req.Visits) {
		return "", fmt.Errorf("count does not match number of visits")
	}

	// Create new job
	jobID := strconv.FormatInt(time.Now().UnixNano(), 10)
	job := models.Job{
		ID:         jobID,
		Status:     models.JobStatusPending,
		VisitCount: req.Count,
		Visits:     req.Visits,
		CreatedAt:  time.Now(),
	}

	// Save job
	if _, err := s.jobRepo.Create(job); err != nil {
		return "", fmt.Errorf("failed to create job: %w", err)
	}

	// Process job asynchronously
	go s.processJob(jobID)

	return jobID, nil
}

func (s *jobService) GetJobStatus(jobID string) (models.JobStatusResponse, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return models.JobStatusResponse{}, err
	}

	response := models.JobStatusResponse{
		Status: job.Status,
		JobID:  job.ID,
	}

	if job.Status == models.JobStatusFailed {
		response.Errors = job.Errors
	}

	return response, nil
}

func (s *jobService) processJob(jobID string) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return
	}

	// Update job status to ongoing
	job.Status = models.JobStatusOngoing
	if err := s.jobRepo.Update(job); err != nil {
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	jobFailed := false

	// Process each visit
	for i, visit := range job.Visits {
		// Verify store exists
		if _, err := s.storeRepo.GetByID(visit.StoreID); err != nil {
			mu.Lock()
			job.Errors = append(job.Errors, models.JobError{
				StoreID: visit.StoreID,
				Error:   "store not found",
			})
			jobFailed = true
			mu.Unlock()
			continue
		}

		// Initialize image results
		job.Visits[i].Results = make([]models.ImageResult, len(visit.ImageURLs))

		// Process each image
		for j, url := range visit.ImageURLs {
			wg.Add(1)
			go func(visitIdx, imgIdx int, imgURL string) {
				defer wg.Done()

				result := models.ImageResult{
					URL:       imgURL,
					Processed: false,
				}

				perimeter, err := s.imageProcessor.ProcessImage(imgURL)
				if err != nil {
					mu.Lock()
					result.Error = err.Error()
					job.Visits[visitIdx].Results[imgIdx] = result
					jobFailed = true
					job.Errors = append(job.Errors, models.JobError{
						StoreID: job.Visits[visitIdx].StoreID,
						Error:   fmt.Sprintf("failed to process image %s: %v", imgURL, err),
					})
					mu.Unlock()
					return
				}

				result.Perimeter = perimeter
				result.Processed = true

				mu.Lock()
				job.Visits[visitIdx].Results[imgIdx] = result
				mu.Unlock()
			}(i, j, url)
		}
	}

	wg.Wait()

	// Update job status
	if jobFailed {
		job.Status = models.JobStatusFailed
	} else {
		job.Status = models.JobStatusCompleted
	}

	completedAt := time.Now()
	job.CompletedAt = &completedAt

	// Update job in repository
	s.jobRepo.Update(job)
}
