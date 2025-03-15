package main

import "net/http"

func (cfg *apiConfig) handlerResetCounter(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("counter reset to 0"))
}
