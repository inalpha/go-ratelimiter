package main

import (
	"net/http"
)

type Limiter struct {
	store Store
}

func NewLimiter(store Store) *Limiter {
	return &Limiter{store}
}
func (l *Limiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := l.store.Get(r)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
