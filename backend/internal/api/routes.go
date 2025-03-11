package api

import (
	"net/http"
)

func SetupRoutes(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/submit/", handler.SubmitJob)
	mux.HandleFunc("/api/status", handler.GetJobStatus)

	return mux
}
