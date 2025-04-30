package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/muszkin/http-server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	const fileRootPath = "."
	const port = "8080"

	apiConfig := apiConfig{fileserverHits: atomic.Int32{}, dbQueries: dbQueries}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.ServeHTTP)
	serveMux.HandleFunc("POST /admin/reset", apiConfig.reset)
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	serveMux.HandleFunc("POST /api/users", apiConfig.handleCreateUserRequest)
	serveMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileRootPath)))))
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	err = server.ListenAndServe()
	if err != nil {
		return
	}
}
