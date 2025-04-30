package main

import (
	"context"
	"net/http"
	"os"
)

func (cfg *apiConfig) reset(w http.ResponseWriter, _ *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
	}
	if err := cfg.dbQueries.RemoveUsers(context.Background()); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hits reset to 0"))
}
