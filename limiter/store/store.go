package store

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Store interface {
	Get(string) *rate.Limiter
	Save(string, *rate.Limiter)
}

type InMemory struct {
	visitors map[string]*visitor
	lock     sync.RWMutex
	*CleanupOptions
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type CleanupOptions struct {
	Rate      time.Duration
	Threshold time.Duration
}

var DefaultCleanupOptions = &CleanupOptions{time.Minute, 3 * time.Minute}

func NewInMemory(c *CleanupOptions) *InMemory {
	im := &InMemory{
		visitors:       make(map[string]*visitor),
		CleanupOptions: c,
	}
	go im.cleanup()
	return im
}

func (im *InMemory) Get(key string) *rate.Limiter {
	im.lock.RLock()
	v, ok := im.visitors[key]
	im.lock.RUnlock()

	if !ok {
		return nil
	}

	im.lock.Lock()
	v.lastSeen = time.Now()
	im.lock.Unlock()

	return v.limiter
}

func (im *InMemory) Save(key string, limiter *rate.Limiter) {
	im.lock.Lock()
	im.visitors[key] = &visitor{limiter, time.Now()}
	im.lock.Unlock()
}

func (im *InMemory) cleanup() {
	for {
		time.Sleep(im.CleanupOptions.Rate)
		im.lock.Lock()
		for ip, v := range im.visitors {
			if time.Now().Sub(v.lastSeen) > 3*im.CleanupOptions.Threshold {
				delete(im.visitors, ip)
			}
		}
		im.lock.Unlock()
	}
}
