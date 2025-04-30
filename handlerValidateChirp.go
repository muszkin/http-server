package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
		Error       string `json:"error"`
	}
	forbidenWords := []string{"kerfuffle", "sharbert", "fornax"}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding body: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something goes wrong")
		return
	}
	if len(params.Body) > 140 {
		resp := response{
			Valid: false,
			Error: "Chirp is too long",
		}
		respondWithJSON(w, http.StatusBadRequest, resp)
		return
	}
	spliitedBody := strings.Split(params.Body, " ")
	for i, word := range spliitedBody {
		if slices.Contains(forbidenWords, strings.ToLower(word)) {
			spliitedBody[i] = "****"
		}
	}
	resp := response{
		CleanedBody: strings.Join(spliitedBody, " "),
		Valid:       true,
	}
	w.WriteHeader(200)
	data, _ := json.Marshal(resp)
	_, _ = w.Write(data)
}
