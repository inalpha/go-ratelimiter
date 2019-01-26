package main

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Store interface {
	Get(r *http.Request) *rate.Limiter
}

type InMemory struct {
	visitors map[string]*visitor
	lock     sync.RWMutex
}

func NewInMemory() *InMemory {
	im := &InMemory{
		visitors: make(map[string]*visitor),
	}
	go im.cleanup()
	return im
}

func (im *InMemory) Get(r *http.Request) *rate.Limiter {
	im.lock.RLock()
	v, ok := im.visitors[r.RemoteAddr]
	im.lock.RUnlock()

	if !ok {
		return im.add(r.RemoteAddr)
	}

	im.lock.Lock()
	v.lastSeen = time.Now()
	im.lock.Unlock()

	return v.limiter
}

func (im *InMemory) add(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(2, 5)
	im.lock.Lock()
	im.visitors[ip] = &visitor{limiter, time.Now()}
	im.lock.Unlock()
	return limiter
}

func (im *InMemory) cleanup() {
	for {
		time.Sleep(time.Minute)
		println("cleaning")

		im.lock.Lock()
		for ip, v := range im.visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				println(ip, "remove")
				delete(im.visitors, ip)
			}
		}
		im.lock.Unlock()
	}
}
