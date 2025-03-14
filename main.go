package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const rootFilePath = "."
	const port = "8080"

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerCountRequests)
	mux.HandleFunc("/reset", apiCfg.handlerResetCounter)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}

}

func (cfg *apiConfig) handlerCountRequests(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", count)))

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
