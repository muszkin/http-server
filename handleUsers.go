package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/muszkin/http-server/internal/auth"
	"github.com/muszkin/http-server/internal/database"
	"net/http"
	"time"
)

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handleCreateUserRequest(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req createUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if len(req.Password) < 1 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 1 character")
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := cfg.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	expiresIn := time.Duration(3600) * time.Second
	userFromDb, err := cfg.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = auth.CheckPasswordHash(userFromDb.HashedPassword, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(userFromDb.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	randomString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	refreshToken, err := cfg.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     randomString,
		UserID:    userFromDb.ID,
		ExpiresAt: time.Now().Add(time.Duration(24*60) * time.Hour),
	})
	userData := loginResponse{
		ID:           userFromDb.ID,
		CreatedAt:    userFromDb.CreatedAt,
		UpdatedAt:    userFromDb.UpdatedAt,
		Email:        userFromDb.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	}
	respondWithJSON(w, http.StatusOK, userData)
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	type refreshTokenResponse struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "missing refresh token")
		return
	}
	userFromRefreshToken, err := cfg.dbQueries.GetUserFromRefreshToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no refresh token found")
		return
	}
	token, err := auth.MakeJWT(userFromRefreshToken.ID, cfg.jwtSecret, time.Second*3600)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseData := refreshTokenResponse{
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, responseData)
}

func (cfg *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "missing refresh token")
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
