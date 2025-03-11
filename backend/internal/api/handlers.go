// internal/api/handlers.go
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"kirana-club-assignment/backend/internal/models"
	"kirana-club-assignment/backend/internal/repository"
	"kirana-club-assignment/backend/internal/service"
)

type Handler struct {
	jobService service.JobService
}

func NewHandler(jobService service.JobService) *Handler {
	return &Handler{
		jobService: jobService,
	}
}

func (h *Handler) SubmitJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate request
	if req.Count <= 0 || len(req.Visits) == 0 || req.Count != len(req.Visits) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request: count must match number of visits"})
		return
	}

	jobID, err := h.jobService.CreateJob(req)
	if err != nil {
		log.Printf("Error creating job: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.JobResponse{JobID: jobID})
}

func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := r.URL.Query().Get("jobid")
	if jobID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	status, err := h.jobService.GetJobStatus(jobID)
	if err != nil {
		if err == repository.ErrJobNotFound {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{})
		} else {
			log.Printf("Error getting job status: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}
