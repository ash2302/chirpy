package main

import (
	"log"
	"net/http"
)

func main() {
	const rootFilePath = "."
	const port = "8080"

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.HandleFunc("/healthz", handlerReadiness)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath))))

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
