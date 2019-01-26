package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	l := NewLimiter(NewInMemory())
	http.ListenAndServe(":4000", l.Handler(mux))
}
func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
