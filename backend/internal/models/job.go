// internal/models/job.go
package models

import (
	"time"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusOngoing   JobStatus = "ongoing"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type Job struct {
	ID          string     `json:"job_id"`
	Status      JobStatus  `json:"status"`
	VisitCount  int        `json:"count"`
	Visits      []Visit    `json:"visits"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Errors      []JobError `json:"errors,omitempty"`
}

type Visit struct {
	StoreID   string        `json:"store_id"`
	ImageURLs []string      `json:"image_url"`
	VisitTime string        `json:"visit_time"`
	Results   []ImageResult `json:"results,omitempty"`
}

type ImageResult struct {
	URL       string  `json:"url"`
	Perimeter float64 `json:"perimeter"`
	Processed bool    `json:"processed"`
	Error     string  `json:"error,omitempty"`
}

type JobError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

type JobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type JobResponse struct {
	JobID string `json:"job_id"`
}

type JobStatusResponse struct {
	Status JobStatus  `json:"status"`
	JobID  string     `json:"job_id"`
	Errors []JobError `json:"error,omitempty"`
}
