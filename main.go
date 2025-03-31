package main

import (
	"database/sql"
	"github.com/ash2302/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	dbQueries := database.New(db)

	const rootFilePath = "."
	const port = "8080"

	apiCfg := apiConfig{
		dbQueries: dbQueries,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCountRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetCounter)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	log.Printf("Serving files from %s on port: %s\n", rootFilePath, port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}

}

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
}
