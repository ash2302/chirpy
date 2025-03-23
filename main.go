package main

import (
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
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCountRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetCounter)

	log.Printf("Serving files from %s on port: %s\n", rootFilePath, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}

}

type apiConfig struct {
	fileServerHits atomic.Int32
}
