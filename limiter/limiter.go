package limiter

import (
	"net/http"
	"ratelimiter/limiter/store"

	"golang.org/x/time/rate"
)

type Limiter struct {
	store store.Store
	KeyGenerator
}

type KeyGenerator func(r *http.Request) string

var DefaultKeyGenerator = func(r *http.Request) string { return r.RemoteAddr }

func NewLimiter(store store.Store, kg KeyGenerator) *Limiter {
	return &Limiter{store, kg}
}

func (l *Limiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := l.get(r)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (l *Limiter) get(r *http.Request) *rate.Limiter {
	key := l.KeyGenerator(r)
	limiter := l.store.Get(key)
	if limiter != nil {
		return limiter
	}
	limiter = rate.NewLimiter(2, 5)
	l.store.Save(key, limiter)
	return limiter
}
