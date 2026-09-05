package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/run", handleRun)
	mux.HandleFunc("POST /api/queue-run", handleQueueRun)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	addr := ":" + port

	fmt.Printf("server: listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

// withCORS allows the web/ dev server (a different origin/port during
// local development) to call this API. This project has no
// authentication or cookies to protect, so a permissive, allow-all
// policy is a reasonable default rather than a security gap.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
