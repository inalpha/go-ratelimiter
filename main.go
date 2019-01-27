package main

import (
	"net/http"

	"ratelimiter/limiter"
	limiterstore "ratelimiter/limiter/store"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	l := limiter.NewLimiter(
		limiterstore.NewInMemory(limiterstore.DefaultCleanupOptions),
		limiter.DefaultKeyGenerator,
	)
	http.ListenAndServe(":4000", l.Handler(mux))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
