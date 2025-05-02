package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/muszkin/http-server/internal/auth"
	"github.com/muszkin/http-server/internal/database"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

type response struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding body: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something goes wrong")
		return
	}
	params.UserId = userId
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Too many characters in body")
		return
	}
	splitBody := strings.Split(params.Body, " ")
	for i, word := range splitBody {
		if slices.Contains(forbiddenWords, strings.ToLower(word)) {
			splitBody[i] = "****"
		}
	}
	chrip, err := cfg.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   strings.Join(splitBody, " "),
		UserID: params.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, response{chrip.ID, chrip.CreatedAt, chrip.UpdatedAt, chrip.Body, chrip.UserID})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	parsedChirps := make([]response, 0)
	for _, chirp := range chirps {
		parsedChirps = append(parsedChirps, response{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, parsedChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpPathId := r.PathValue("chirpId")
	chirpId, err := uuid.Parse(chirpPathId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	chirp, err := cfg.dbQueries.GetChrip(context.Background(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}
	chirpData := response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirpData)
}
